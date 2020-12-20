package command

import (
	commonCommand "github.com/klwxsrx/expense-tracker/pkg/common/app/command"
	"github.com/klwxsrx/expense-tracker/pkg/expense/domain"
)

var ResultMap = map[error]commonCommand.Result{
	domain.ErrorInvalidAccountTitle:      commonCommand.ResultInvalidArgument,
	domain.ErrorAlreadyDeletedAccount:    commonCommand.ResultNotFound,
	domain.ErrorAccountTitleIsDuplicated: commonCommand.ResultDuplicateConflict,
	domain.ErrorAccountIsNotExists:       commonCommand.ResultNotFound,
	domain.ErrorInvalidCurrency:          commonCommand.ResultInvalidArgument,
}
