package mysql

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
)

type Client interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Get(dest interface{}, query string, args ...interface{}) error
	Select(dest interface{}, query string, args ...interface{}) error
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
