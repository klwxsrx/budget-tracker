package mysql

import (
	"errors"
)

const lockTimeout = 10

type Lock interface {
	Get() error
	Release()
}

type lock struct {
	client Client
	name   string
}

func (l *lock) Get() error {
	var success int
	err := l.client.Get(&success, "SELECT GET_LOCK(MD5(CONCAT(DATABASE(), ?)), ?)", l.name, lockTimeout)
	if err != nil {
		return err
	}
	if success == 0 {
		return errors.New("lock attempt timed out")
	}
	return nil
}

func (l *lock) Release() {
	_ = l.client.Get("SELECT RELEASE_LOCK(MD5(CONCAT(DATABASE(), ?)))", l.name)
}

func NewLock(client Client, name string) Lock {
	return &lock{client, name}
}
