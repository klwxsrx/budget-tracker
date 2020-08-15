package account

type Specification interface {
	IsSatisfied(a *Account) bool
}

type titleSpecification struct {
	title string
}

func (s *titleSpecification) IsSatisfied(a *Account) bool {
	return a.state.Title == s.title
}

type Repository interface {
	NextID() ID
	Update(a *Account) error
	GetByID(id ID) (*Account, error)
	Exists(spec Specification) (bool, error)
}
