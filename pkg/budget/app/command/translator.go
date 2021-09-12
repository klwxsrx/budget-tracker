package command

import (
	"errors"

	"github.com/klwxsrx/budget-tracker/pkg/budget/domain"
	commonappcommand "github.com/klwxsrx/budget-tracker/pkg/common/app/command"
)

type translator struct {
	resultMap map[error]commonappcommand.Result
}

func (t *translator) Translate(err error) commonappcommand.Result {
	if err == nil {
		return commonappcommand.ResultSuccess
	}

	for e, r := range t.resultMap {
		if errors.Is(err, e) {
			return r
		}
	}
	return commonappcommand.ResultUnknownError
}

func NewTranslator() commonappcommand.ErrorTranslator {
	var resultMap = map[error]commonappcommand.Result{
		domain.ErrAccountListAlreadyExists: commonappcommand.ResultDuplicateConflict,
		domain.ErrAccountListDoesNotExist:  commonappcommand.ResultNotFound,
		domain.ErrAccountDuplicateTitle:    commonappcommand.ResultDuplicateConflict,
		domain.ErrAccountDoesNotExist:      commonappcommand.ResultNotFound,
		domain.ErrAccountInvalidTitle:      commonappcommand.ResultInvalidArgument,
	}
	return &translator{resultMap}
}
