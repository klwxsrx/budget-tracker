package main

import (
	"fmt"
	"github.com/klwxsrx/expense-tracker/pkg/common/infrastructure/mysql"
	"github.com/klwxsrx/expense-tracker/pkg/expense/app/command"
	"github.com/klwxsrx/expense-tracker/pkg/expense/infrastructure"
	"os"
)

func main() {
	db := mysql.NewDatabase(mysql.Dsn{
		User:     "expense",
		Password: "1234",
		Address:  "expense-db:3306",
		Database: "expense",
	}, 5)

	err := db.OpenConnection()
	if err != nil {
		panic(err)
	}
	client, err := db.GetClient()
	if err != nil {
		panic(err)
	}
	migrationDir, ok := os.LookupEnv("EVENT_STORE_MIGRATIONS_DIR")
	if !ok {
		panic("cannot find migrations dir")
	}
	migration, err := mysql.NewMigration(client, migrationDir)
	if err != nil {
		panic(err)
	}
	err = migration.Migrate()
	if err != nil {
		panic(err)
	}

	bus := infrastructure.NewContainer(client).CommandBus()
	cmd := command.CreateAccount{Title: "Some", Currency: "USD", InitialBalance: 42}
	err = bus.Publish(&cmd)
	fmt.Println(err)
}
