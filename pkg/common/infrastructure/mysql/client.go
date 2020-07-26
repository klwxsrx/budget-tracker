package mysql

import (
	"github.com/jmoiron/sqlx"
)

type Client interface {
	Get(dest interface{}, query string, args ...interface{}) error
	// TODO: add required methods from sqlx
}

type ClientTransaction interface {
	Client
	Commit() error
	Rollback() error
}

type TransactionalClient interface {
	Client
	Begin() (ClientTransaction, error)
}

type client struct {
	*sqlx.DB
}

func (c *client) Begin() (ClientTransaction, error) {
	return c.Beginx()
}
