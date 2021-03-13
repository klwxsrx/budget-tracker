package domain

type AccountSpecification interface {
	IsSatisfied(a *Account) bool
}

type accountTitleSpecification struct {
	title string
}

func (s *accountTitleSpecification) IsSatisfied(a *Account) bool {
	return a.state.Title == s.title
}

type AccountRepository interface {
	NextID() AccountID
	Update(a *Account) error
	GetByID(id AccountID) (*Account, error)
	Exists(spec AccountSpecification) (bool, error)
}
