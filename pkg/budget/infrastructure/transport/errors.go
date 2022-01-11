package transport

import (
	"errors"
	"net/http"

	"github.com/klwxsrx/budget-tracker/pkg/budget/domain"
)

// nolint:gochecknoglobals
var errorsToHTTPCode = map[error]int{
	domain.ErrAccountListAlreadyExists: http.StatusConflict,
	domain.ErrAccountListDoesNotExist:  http.StatusNotFound,
	domain.ErrAccountDoesNotExist:      http.StatusNotFound,
	domain.ErrAccountInvalidTitle:      http.StatusBadRequest,
	domain.ErrAccountDuplicateTitle:    http.StatusConflict,
}

func translateError(err error) (errorCode int) {
	if err == nil {
		return http.StatusNoContent
	}
	for e, code := range errorsToHTTPCode {
		if errors.Is(err, e) {
			return code
		}
	}
	return http.StatusInternalServerError
}
