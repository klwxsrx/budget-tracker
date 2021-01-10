package mysql

import (
	"fmt"
	"github.com/cenkalti/backoff"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/klwxsrx/expense-tracker/pkg/common/app/logger"
	"time"
)

const (
	maxConnectionTime  = time.Minute
	maxOpenConnections = 10
)

type Dsn struct {
	User     string
	Password string
	Address  string
	Database string
}

func (d Dsn) String() string {
	return fmt.Sprintf("%v:%v@(%v)/%v", d.User, d.Password, d.Address, d.Database)
}

type Connection interface {
	Client() (TransactionalClient, error)
	Close()
}

type connection struct {
	config Dsn
	db     *sqlx.DB
	logger logger.Logger
}

func (c *connection) Client() (TransactionalClient, error) {
	return &client{c.db}, nil
}

func (c *connection) Close() {
	err := c.db.Close()
	if err != nil {
		c.logger.WithError(err).Error("failed to close mongo db connection")
	}
}

func (c *connection) openConnection() error {
	var err error
	c.db, err = sqlx.Open("mysql", c.config.String()+"?parseTime=true")
	c.db.SetMaxOpenConns(maxOpenConnections)
	err = backoff.Retry(func() error {
		return c.db.Ping()
	}, newOpenConnectionBackoff())
	if err != nil {
		_ = c.db.Close()
		return fmt.Errorf("failed to open mysql connection: %v", err)
	}
	return nil
}

func newOpenConnectionBackoff() *backoff.ExponentialBackOff {
	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = maxConnectionTime
	return b
}

func NewConnection(config Dsn, logger logger.Logger) (Connection, error) {
	db := &connection{config: config, logger: logger}
	err := db.openConnection()
	return db, err
}
