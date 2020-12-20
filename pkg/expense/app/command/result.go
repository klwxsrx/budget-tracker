package command

import (
	commonCommand "github.com/klwxsrx/expense-tracker/pkg/common/app/command"
	"github.com/klwxsrx/expense-tracker/pkg/expense/domain"
)

var ResultMap = map[error]commonCommand.Result{
	domain.InvalidAccountTitleError:      commonCommand.ResultInvalidArgument,
	domain.AlreadyDeletedAccountError:    commonCommand.ResultNotFound,
	domain.AccountTitleIsDuplicatedError: commonCommand.ResultDuplicateConflict,
	domain.AccountIsNotExistsError:       commonCommand.ResultNotFound,
	domain.InvalidCurrencyError:          commonCommand.ResultInvalidArgument,
}
