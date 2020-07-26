package mysql

import (
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
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

type Database interface {
	OpenConnection() error
	CloseConnection() error
	GetClient() (TransactionalClient, error)
}

type database struct {
	config Dsn
	db     *sqlx.DB
}

func (d *database) OpenConnection() error {
	var err error
	if d.db != nil {
		err = d.db.Close()
	}
	if err != nil {
		return err
	}

	d.db, err = sqlx.Open("mysql", d.config.String())
	if err == nil {
		// Execute query to check connectivity. Ping method isn't working for some drivers
		_, err = d.db.Exec("SELECT TRUE AS test_connection")
	}
	return err
}

func (d *database) CloseConnection() error {
	if d.db == nil {
		return errors.New("connection is already closed")
	}
	return d.db.Close()
}

func (d *database) GetClient() (TransactionalClient, error) {
	if d.db == nil {
		return nil, errors.New("connection is not opened")
	}
	return &client{d.db}, nil
}

func NewDatabase(config Dsn) Database {
	return &database{config: config}
}
