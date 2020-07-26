package domain

type AccountRepository interface {
	Update(a *Account) error
	GetByID(id AccountID) (*Account, error)
	Exists(title string) (bool, error)
}
