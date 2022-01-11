package transport

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"github.com/klwxsrx/budget-tracker/pkg/budget/app/command"
	commonappcommand "github.com/klwxsrx/budget-tracker/pkg/common/app/command"
	"github.com/klwxsrx/budget-tracker/pkg/common/app/logger"
)

var (
	errInvalidParameter = errors.New("invalid parameter")
	errEmptyJSONBody    = errors.New("empty json body")
)

type commandParser func(r *http.Request) (commonappcommand.Command, error)

type route struct {
	Name    string
	Method  string
	Pattern string
	Parser  commandParser
}

func getRoutes() []route {
	return []route{
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
}

type addAccountBody struct {
	Title          string `json:"title"`
	InitialBalance int    `json:"initialBalance"`
}

func addAccountParser(r *http.Request) (commonappcommand.Command, error) {
	budgetID, err := parseUUID(mux.Vars(r)["budgetID"])
	if err != nil {
		return nil, errInvalidParameter
	}
	var body addAccountBody
	if err := parseJSONFromBody(r, &body); err != nil {
		return nil, err
	}
	if body.Title == "" {
		return nil, errInvalidParameter
	}
	return command.NewAccountAdd(budgetID, body.Title, body.InitialBalance), nil
}

func reorderAccountParser(r *http.Request) (commonappcommand.Command, error) {
	budgetID, err := parseUUID(mux.Vars(r)["budgetID"])
	if err != nil {
		return nil, errInvalidParameter
	}
	accountID, err := parseUUID(mux.Vars(r)["accountID"])
	if err != nil {
		return nil, errInvalidParameter
	}
	position, err := parseInt(mux.Vars(r)["position"])
	if err != nil {
		return nil, errInvalidParameter
	}
	return command.NewAccountReorder(budgetID, accountID, position), nil
}

type renameAccountBody struct {
	Title string `json:"title"`
}

func renameAccountParser(r *http.Request) (commonappcommand.Command, error) {
	budgetID, err := parseUUID(mux.Vars(r)["budgetID"])
	if err != nil {
		return nil, errInvalidParameter
	}
	accountID, err := parseUUID(mux.Vars(r)["accountID"])
	if err != nil {
		return nil, errInvalidParameter
	}
	var body renameAccountBody
	if err := parseJSONFromBody(r, &body); err != nil {
		return nil, err
	}
	if body.Title == "" {
		return nil, errInvalidParameter
	}
	return command.NewAccountRename(budgetID, accountID, body.Title), nil
}

const (
	AccountStatusActive    = "active"
	AccountStatusCancelled = "cancelled"
)

func changeAccountStatusParser(r *http.Request) (commonappcommand.Command, error) {
	budgetID, err := parseUUID(mux.Vars(r)["budgetID"])
	if err != nil {
		return nil, errInvalidParameter
	}
	accountID, err := parseUUID(mux.Vars(r)["accountID"])
	if err != nil {
		return nil, errInvalidParameter
	}
	status, ok := mux.Vars(r)["status"]
	if !ok {
		return nil, errInvalidParameter
	}
	switch status {
	case AccountStatusActive:
		return command.NewAccountActivate(budgetID, accountID), nil
	case AccountStatusCancelled:
		return command.NewAccountCancel(budgetID, accountID), nil
	default:
		return nil, errInvalidParameter
	}
}

func deleteAccountParser(r *http.Request) (commonappcommand.Command, error) {
	budgetID, err := parseUUID(mux.Vars(r)["budgetID"])
	if err != nil {
		return nil, errInvalidParameter
	}
	accountID, err := parseUUID(mux.Vars(r)["accountID"])
	if err != nil {
		return nil, errInvalidParameter
	}
	return command.NewAccountDelete(budgetID, accountID), nil
}

func parseUUID(str string) (uuid.UUID, error) {
	return uuid.Parse(str)
}

func parseInt(str string) (int, error) {
	return strconv.Atoi(str)
}

func parseJSONFromBody(r *http.Request, v interface{}) error {
	err := json.NewDecoder(r.Body).Decode(v)
	if errors.Is(err, io.EOF) {
		return errEmptyJSONBody
	}
	return err
}

func getHandlerFunc(bus commonappcommand.Bus, parser commandParser) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cmd, err := parser(r)
		switch {
		case errors.Is(err, errInvalidParameter):
		case errors.Is(err, errEmptyJSONBody):
			w.WriteHeader(http.StatusBadRequest)
			return
		case err != nil:
			w.WriteHeader(http.StatusInternalServerError)
			return
		default:
		}

		err = bus.Publish(cmd)

		httpCode := translateError(err)
		w.WriteHeader(httpCode)
	}
}

func NewHTTPHandler(bus commonappcommand.Bus, loggerImpl logger.Logger) http.Handler {
	r := mux.NewRouter().PathPrefix("/budget").Subrouter()
	for _, route := range getRoutes() {
		r.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			HandlerFunc(getHandlerFunc(bus, route.Parser))
	}
	r.Use(getLoggingMiddleware(loggerImpl))
	return r
}
