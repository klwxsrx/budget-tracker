package service

type UnitOfWork interface {
	Execute(f func(r RepositoryProvider) error) error
}

type RepositoryProvider interface{}
