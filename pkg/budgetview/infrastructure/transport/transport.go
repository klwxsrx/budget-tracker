package transport

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"github.com/klwxsrx/budget-tracker/pkg/budgetview/app/model"
	"github.com/klwxsrx/budget-tracker/pkg/budgetview/infrastructure"
	"github.com/klwxsrx/budget-tracker/pkg/common/app/logger"
	"github.com/klwxsrx/budget-tracker/pkg/common/infrastructure/transport"
)

type route struct {
	Name    string
	Method  string
	Pattern string
	Handler func(writer http.ResponseWriter, request *http.Request, container infrastructure.Container)
}

func getRoutes() []route {
	return []route{
		{
			"ListBudgets",
			http.MethodGet,
			"/budget/list",
			listBudgetsHandler,
		},
		{
			"ListAccounts",
			http.MethodGet,
			"/account/{budgetID}/list",
			listAccountsHandler,
		},
	}
}

type BudgetJSON struct {
	ID       uuid.UUID `json:"id"`
	Title    string    `json:"title"`
	Currency string    `json:"currency"`
}

func listBudgetsHandler(
	w http.ResponseWriter,
	_ *http.Request,
	container infrastructure.Container,
) {
	budgets, err := container.BudgetQueryService().ListBudgets()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	result := make([]BudgetJSON, 0, len(budgets))
	for _, budget := range budgets {
		result = append(result, BudgetJSON{
			ID:       budget.ID,
			Title:    budget.Title,
			Currency: budget.Currency,
		})
	}

	writeJSONResult(w, result)
}

type AccountJSON struct {
	ID             uuid.UUID `json:"id"`
	BudgetID       uuid.UUID `json:"budgetID"`
	Title          string    `json:"title"`
	Status         string    `json:"status"`
	InitialBalance int       `json:"initialBalance"`
	CurrentBalance int       `json:"currentBalance"`
}

func listAccountsHandler(
	w http.ResponseWriter,
	r *http.Request,
	container infrastructure.Container,
) {
	budgetID, err := parseUUID(mux.Vars(r)["budgetID"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	exists, err := container.BudgetQueryService().ExistByIDs([]uuid.UUID{budgetID})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	accounts, err := container.AccountQueryService().ListAccounts(budgetID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	result := make([]AccountJSON, 0, len(accounts))
	for _, account := range accounts {
		result = append(result, AccountJSON{
			ID:             account.ID,
			BudgetID:       account.BudgetID,
			Title:          account.Title,
			Status:         toStringAccountStatus(account.Status),
			InitialBalance: account.InitialBalance,
			CurrentBalance: account.CurrentBalance,
		})
	}

	writeJSONResult(w, result)
}

func writeJSONResult(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func parseUUID(str string) (uuid.UUID, error) {
	return uuid.Parse(str)
}

func toStringAccountStatus(status int) string {
	switch status {
	case model.AccountStatusActive:
		return "active"
	case model.AccountStatusCancelled:
	}
	return "cancelled"
}

func getHandlerFunc(
	container infrastructure.Container,
	handler func(writer http.ResponseWriter, request *http.Request, container infrastructure.Container),
) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		handler(writer, request, container)
	}
}

func addLivenessCheckRoute(router *mux.Router) {
	router.
		Methods(http.MethodGet).
		Path("/internal/status/live").
		HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("ok"))
		})
}

func NewHTTPHandler(container infrastructure.Container, loggerImpl logger.Logger) http.Handler {
	router := mux.NewRouter()
	addLivenessCheckRoute(router)

	api := router.PathPrefix("/budget-view").Subrouter()
	for _, route := range getRoutes() {
		api.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			HandlerFunc(getHandlerFunc(container, route.Handler))
	}
	router.Use(transport.GetLoggingMiddleware(loggerImpl))
	return router
}
