package account

import "github.com/klwxsrx/expense-tracker/pkg/common/app/command"

const (
	createCommandType = "expense.account.create"
	renameCommandType = "expense.account.rename"
	deleteCommandType = "expense.account.delete"
)

type CreateCommand struct {
	Title          string
	Currency       string
	InitialBalance int
}

func (c *CreateCommand) GetType() command.Type {
	return createCommandType
}

type RenameCommand struct {
	ID    string
	Title string
}

func (c *RenameCommand) GetType() command.Type {
	return renameCommandType
}

type DeleteCommand struct {
	ID string
}

func (c *DeleteCommand) GetType() command.Type {
	return deleteCommandType
}
