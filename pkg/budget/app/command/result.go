package command

import (
	"github.com/klwxsrx/budget-tracker/pkg/budget/domain"
	commonCommand "github.com/klwxsrx/budget-tracker/pkg/common/app/command"
)

var ResultMap = map[error]commonCommand.Result{
	domain.ErrorInvalidAccountTitle:      commonCommand.ResultInvalidArgument,
	domain.ErrorAlreadyDeletedAccount:    commonCommand.ResultNotFound,
	domain.ErrorAccountTitleIsDuplicated: commonCommand.ResultDuplicateConflict,
	domain.ErrorAccountIsNotExists:       commonCommand.ResultNotFound,
	domain.ErrorInvalidCurrency:          commonCommand.ResultInvalidArgument,
}
