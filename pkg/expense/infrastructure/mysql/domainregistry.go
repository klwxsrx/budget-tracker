package mysql

import (
	eventApp "github.com/klwxsrx/expense-tracker/pkg/common/app/event"
	eventMysql "github.com/klwxsrx/expense-tracker/pkg/common/infrastructure/event/mysql"
	"github.com/klwxsrx/expense-tracker/pkg/common/infrastructure/mysql"
	"github.com/klwxsrx/expense-tracker/pkg/expense/app/command"
	"github.com/klwxsrx/expense-tracker/pkg/expense/domain"
	"github.com/klwxsrx/expense-tracker/pkg/expense/infrastructure/serialization"
)

type domainRegistry struct {
	client       mysql.Client
	dispatcher   eventApp.Dispatcher
	serializer   eventApp.Serializer
	deserializer serialization.EventDeserializer
}

func (dr *domainRegistry) AccountService() domain.AccountService {
	store := eventMysql.NewStore(dr.client, dr.serializer)
	repo := NewAccountRepository(dr.dispatcher, store, dr.deserializer)
	return domain.NewAccountService(repo)
}

func newDomainRegistry(
	client mysql.Client,
	dispatcher eventApp.Dispatcher,
	serializer eventApp.Serializer,
	deserializer serialization.EventDeserializer,
) command.DomainRegistry {
	return &domainRegistry{client, dispatcher, serializer, deserializer}
}
