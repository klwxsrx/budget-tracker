package transport

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/klwxsrx/budget-tracker/pkg/budget/app/command"
	appCommand "github.com/klwxsrx/budget-tracker/pkg/common/app/command"
	"github.com/klwxsrx/budget-tracker/pkg/common/app/logger"
	"net/http"
)

var errorInvalidParameter = errors.New("invalid parameter")

type commandParser func(r *http.Request) (appCommand.Command, error)

type route struct {
	Name    string
	Method  string
	Pattern string
	Parser  commandParser
}

var routes = []route{
	{
		"CreateAccount",
		"POST",
		"/account",
		createAccountParser,
	},
	{
		"RenameAccount",
		"PUT",
		"/account/{accountId}/title",
		renameAccountParser,
	},
	{
		"DeleteAccount",
		"DELETE",
		"/account/{accountId}",
		deleteAccountParser,
	},
}

type createAccountBody struct {
	Title          string `json:"title"`
	Currency       string `json:"currency"`
	InitialBalance int    `json:"initialBalance"`
}

type renameAccountBody struct {
	Title string `json:"title"`
}

func createAccountParser(r *http.Request) (appCommand.Command, error) {
	var body createAccountBody
	if err := parseJsonFromBody(r, &body); err != nil {
		return nil, err
	}
	if body.Title == "" || body.Currency == "" {
		return nil, errorInvalidParameter
	}
	return &command.CreateAccount{Title: body.Title, Currency: body.Currency, InitialBalance: body.InitialBalance}, nil
}

func renameAccountParser(r *http.Request) (appCommand.Command, error) {
	var body renameAccountBody
	if err := parseJsonFromBody(r, &body); err != nil {
		return nil, err
	}

	accountID, err := parseUuid(mux.Vars(r)["accountId"])
	if body.Title == "" || err != nil {
		return nil, errorInvalidParameter
	}
	return &command.RenameAccount{ID: accountID, Title: body.Title}, nil
}

func deleteAccountParser(r *http.Request) (appCommand.Command, error) {
	accountID, err := parseUuid(mux.Vars(r)["accountId"])
	if err != nil {
		return nil, err
	}
	return &command.DeleteAccount{ID: accountID}, nil
}

func parseUuid(str string) (uuid.UUID, error) {
	return uuid.Parse(str)
}

func parseJsonFromBody(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func getHandlerFunc(bus appCommand.Bus, parser commandParser) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cmd, err := parser(r)
		switch {
		case errors.Is(err, errorInvalidParameter):
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
			w.WriteHeader(http.StatusInternalServerError)
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
