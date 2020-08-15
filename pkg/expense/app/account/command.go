package account

import "github.com/klwxsrx/expense-tracker/pkg/common/app/command"

const (
	createAccountType = "expense.create_account"
	renameAccountType = "expense.rename_account"
	deleteAccountType = "expense.delete_account"
)

type CreateAccount struct {
	Title          string
	Currency       string
	InitialBalance int
}

func (c *CreateAccount) GetType() command.Type {
	return createAccountType
}

type RenameAccount struct {
	ID    string
	Title string
}

func (c *RenameAccount) GetType() command.Type {
	return renameAccountType
}

type DeleteAccount struct {
	ID string
}

func (c *DeleteAccount) GetType() command.Type {
	return deleteAccountType
}
