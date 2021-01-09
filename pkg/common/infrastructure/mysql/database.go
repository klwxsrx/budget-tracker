package mysql

import (
	"errors"
	"fmt"
	"github.com/cenkalti/backoff"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"time"
)

const maxConnectionTime = time.Minute

type Dsn struct {
	User     string
	Password string
	Address  string
	Database string
}

func (d Dsn) String() string {
	return fmt.Sprintf("%v:%v@(%v)/%v", d.User, d.Password, d.Address, d.Database)
}

type Database interface {
	Client() (TransactionalClient, error)
	Close()
}

type database struct {
	config  Dsn
	maxConn int
	db      *sqlx.DB
}

func (d *database) Client() (TransactionalClient, error) {
	if d.db == nil {
		return nil, errors.New("connection is not opened")
	}
	return &client{d.db}, nil
}

func (d *database) Close() {
	_ = d.db.Close()
}

func (d *database) openConnection() error {
	var err error
	d.db, err = sqlx.Open("mysql", d.config.String()+"?parseTime=true")
	d.db.SetMaxOpenConns(d.maxConn)
	err = backoff.Retry(func() error {
		return d.db.Ping()
	}, newOpenConnectionBackoff())
	if err != nil {
		_ = d.db.Close()
		return fmt.Errorf("failed to open the connection: %v", err)
	}
	return nil
}

func newOpenConnectionBackoff() *backoff.ExponentialBackOff {
	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = maxConnectionTime
	return b
}

func NewDatabase(config Dsn, maxConnections int) (Database, error) {
	db := &database{config: config, maxConn: maxConnections}
	err := db.openConnection()
	return db, err
}
