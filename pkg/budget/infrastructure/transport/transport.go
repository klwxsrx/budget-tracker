package transport

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/klwxsrx/budget-tracker/pkg/budget/app/command"
	appCommand "github.com/klwxsrx/budget-tracker/pkg/common/app/command"
	"github.com/klwxsrx/budget-tracker/pkg/common/app/logger"
	"io"
	"net/http"
	"strconv"
)

var errorInvalidParameter = errors.New("invalid parameter")
var errorEmptyJsonBody = errors.New("empty json body")

type commandParser func(r *http.Request) (appCommand.Command, error)

type route struct {
	Name    string
	Method  string
	Pattern string
	Parser  commandParser
}

var routes = []route{
	{
		"AddAccount",
		"POST",
		"/account/{budgetID}",
		addAccountParser,
	},
	{
		"ReorderAccount",
		"PUT",
		"/account/{budgetID}/{accountID}/order/{position}",
		reorderAccountParser,
	},
	{
		"RenameAccount",
		"PUT",
		"/account/{budgetID}/{accountID}/title",
		renameAccountParser,
	},
	{
		"ChangeAccountStatus",
		"PUT",
		"/account/{budgetID}/{accountID}/status/{status}",
		changeAccountStatusParser,
	},
	{
		"DeleteAccount",
		"DELETE",
		"/account/{budgetID}/{accountID}",
		deleteAccountParser,
	},
}

type addAccountBody struct {
	Title          string `json:"title"`
	Currency       string `json:"currency"`
	InitialBalance int    `json:"initialBalance"`
}

func addAccountParser(r *http.Request) (appCommand.Command, error) {
	budgetID, err := parseUuid(mux.Vars(r)["budgetID"])
	if err != nil {
		return nil, errorInvalidParameter
	}
	var body addAccountBody
	if err := parseJsonFromBody(r, &body); err != nil {
		return nil, err
	}
	if body.Title == "" || body.Currency == "" {
		return nil, errorInvalidParameter
	}
	return command.NewAccountAdd(budgetID, body.Title, body.Currency, body.InitialBalance), nil
}

func reorderAccountParser(r *http.Request) (appCommand.Command, error) {
	budgetID, err := parseUuid(mux.Vars(r)["budgetID"])
	if err != nil {
		return nil, errorInvalidParameter
	}
	accountID, err := parseUuid(mux.Vars(r)["accountID"])
	if err != nil {
		return nil, errorInvalidParameter
	}
	position, err := parseInt(mux.Vars(r)["position"])
	if err != nil {
		return nil, errorInvalidParameter
	}
	return command.NewAccountReorder(budgetID, accountID, position), nil
}

type renameAccountBody struct {
	Title string `json:"title"`
}

func renameAccountParser(r *http.Request) (appCommand.Command, error) {
	budgetID, err := parseUuid(mux.Vars(r)["budgetID"])
	if err != nil {
		return nil, errorInvalidParameter
	}
	accountID, err := parseUuid(mux.Vars(r)["accountID"])
	if err != nil {
		return nil, errorInvalidParameter
	}
	var body renameAccountBody
	if err := parseJsonFromBody(r, &body); err != nil {
		return nil, err
	}
	if body.Title == "" {
		return nil, errorInvalidParameter
	}
	return command.NewAccountRename(budgetID, accountID, body.Title), nil
}

const (
	AccountStatusActive    = "active"
	AccountStatusCancelled = "cancelled"
)

func changeAccountStatusParser(r *http.Request) (appCommand.Command, error) {
	budgetID, err := parseUuid(mux.Vars(r)["budgetID"])
	if err != nil {
		return nil, errorInvalidParameter
	}
	accountID, err := parseUuid(mux.Vars(r)["accountID"])
	if err != nil {
		return nil, errorInvalidParameter
	}
	status, ok := mux.Vars(r)["status"]
	if !ok {
		return nil, errorInvalidParameter
	}
	switch status {
	case AccountStatusActive:
		return command.NewAccountActivate(budgetID, accountID), nil
	case AccountStatusCancelled:
		return command.NewAccountCancel(budgetID, accountID), nil
	default:
		return nil, errorInvalidParameter
	}
}

func deleteAccountParser(r *http.Request) (appCommand.Command, error) {
	budgetID, err := parseUuid(mux.Vars(r)["budgetID"])
	if err != nil {
		return nil, errorInvalidParameter
	}
	accountID, err := parseUuid(mux.Vars(r)["accountID"])
	if err != nil {
		return nil, errorInvalidParameter
	}
	return command.NewAccountDelete(budgetID, accountID), nil
}

func parseUuid(str string) (uuid.UUID, error) {
	return uuid.Parse(str)
}

func parseInt(str string) (int, error) {
	return strconv.Atoi(str)
}

func parseJsonFromBody(r *http.Request, v interface{}) error {
	err := json.NewDecoder(r.Body).Decode(v)
	if errors.Is(err, io.EOF) {
		return errorEmptyJsonBody
	}
	return err
}

func getHandlerFunc(bus appCommand.Bus, parser commandParser) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cmd, err := parser(r)
		switch {
		case errors.Is(err, errorInvalidParameter):
		case errors.Is(err, errorEmptyJsonBody):
			w.WriteHeader(http.StatusBadRequest)
			return
		case err != nil:
			w.WriteHeader(http.StatusInternalServerError)
			return
		default:
		}

		switch bus.Publish(cmd) {
		case appCommand.ResultSuccess:
			w.WriteHeader(http.StatusNoContent)
		case appCommand.ResultInvalidArgument:
			w.WriteHeader(http.StatusBadRequest)
		case appCommand.ResultNotFound:
			w.WriteHeader(http.StatusNotFound)
		case appCommand.ResultDuplicateConflict:
			w.WriteHeader(http.StatusConflict)
		case appCommand.ResultUnknownError:
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func NewHttpHandler(bus appCommand.Bus, logger logger.Logger) http.Handler {
	r := mux.NewRouter()
	for _, route := range routes {
		r.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			HandlerFunc(getHandlerFunc(bus, route.Parser))
	}
	r.Use(getLoggingMiddleware(logger))
	return r
}
