package domain

type AccountSpecification interface {
	IsSatisfied(a *Account) bool
}

type accountWithTitleSpecification struct {
	title string
}

func (s *accountWithTitleSpecification) IsSatisfied(a *Account) bool {
	return a.state.Title == s.title
}

type AccountRepository interface {
	Update(a *Account) error
	GetByID(id AccountID) (*Account, error)
	Exists(spec AccountSpecification) (bool, error)
}
