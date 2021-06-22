package command

import (
	"github.com/klwxsrx/budget-tracker/pkg/budget/domain"
	commonCommand "github.com/klwxsrx/budget-tracker/pkg/common/app/command"
)

var ResultMap = map[error]commonCommand.Result{
	domain.ErrorCurrencyInvalid:          commonCommand.ResultInvalidArgument,
	domain.ErrorAccountListAlreadyExists: commonCommand.ResultDuplicateConflict,
	domain.ErrorAccountListDoesNotExist:  commonCommand.ResultNotFound,
	domain.ErrorAccountDuplicateTitle:    commonCommand.ResultDuplicateConflict,
	domain.ErrorAccountDoesNotExist:      commonCommand.ResultNotFound,
	domain.ErrorAccountInvalidTitle:      commonCommand.ResultInvalidArgument,
}
