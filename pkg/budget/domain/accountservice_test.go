package domain

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"testing"
)

type mockRepo struct {
	items map[BudgetID]*AccountList
}

func (repo *mockRepo) FindByID(id BudgetID) (*AccountList, error) {
	return repo.items[id], nil
}

func (repo *mockRepo) Update(list *AccountList) error {
	repo.items[list.GetID()] = list
	return nil
}

func TestAccountService_Create(t *testing.T) {
	id := BudgetID{uuid.New()}
	repo := &mockRepo{make(map[BudgetID]*AccountList)}
	service := NewAccountListService(repo)

	err := service.Create(id)
	assert.NoError(t, err)
	list, err := repo.FindByID(id)
	assert.NoError(t, err)

	assert.Len(t, list.GetChanges(), 1)
	listCreatedEvent, ok := list.GetChanges()[0].(*AccountListCreatedEvent)
	require.True(t, ok)
	assert.Equal(t, id.String(), listCreatedEvent.AggregateID().String())
	assert.Equal(t, accountListAggregateName, listCreatedEvent.AggregateName())
	assert.Equal(t, EventTypeAccountListCreated, listCreatedEvent.Type())

	assert.Equal(t, id.String(), list.GetID().String())
	assert.Len(t, list.getAccounts(), 0)

	err = service.Create(id)
	assert.EqualError(t, err, ErrAccountListAlreadyExists.Error())
}

func TestAccountService_Add(t *testing.T) {
	budgetID := BudgetID{uuid.New()}
	repo := &mockRepo{make(map[BudgetID]*AccountList)}
	service := NewAccountListService(repo)
	err := service.Create(budgetID)
	assert.NoError(t, err)

	_, err = service.Add(BudgetID{uuid.New()}, "some", MoneyAmount(4200))
	assert.EqualError(t, err, ErrAccountListDoesNotExist.Error())

	_, err = service.Add(budgetID, "new", MoneyAmount(1300))
	assert.NoError(t, err)
}

func TestAccountService_Reorder(t *testing.T) {
	budgetID := BudgetID{uuid.New()}
	repo := &mockRepo{make(map[BudgetID]*AccountList)}
	service := NewAccountListService(repo)
	err := service.Create(budgetID)
	assert.NoError(t, err)

	err = service.Reorder(BudgetID{uuid.New()}, AccountID{uuid.New()}, 0)
	assert.EqualError(t, err, ErrAccountListDoesNotExist.Error())

	id, err := service.Add(budgetID, "new", MoneyAmount(1300))
	assert.NoError(t, err)
	err = service.Reorder(budgetID, id, 0)
	assert.NoError(t, err)
}

func TestAccountService_Rename(t *testing.T) {
	budgetID := BudgetID{uuid.New()}
	repo := &mockRepo{make(map[BudgetID]*AccountList)}
	service := NewAccountListService(repo)
	err := service.Create(budgetID)
	assert.NoError(t, err)

	err = service.Rename(BudgetID{uuid.New()}, AccountID{uuid.New()}, "some")
	assert.EqualError(t, err, ErrAccountListDoesNotExist.Error())

	id, err := service.Add(budgetID, "new", MoneyAmount(1300))
	assert.NoError(t, err)
	err = service.Rename(budgetID, id, "new")
	assert.NoError(t, err)
}

func TestAccountService_Activate(t *testing.T) {
	budgetID := BudgetID{uuid.New()}
	repo := &mockRepo{make(map[BudgetID]*AccountList)}
	service := NewAccountListService(repo)
	err := service.Create(budgetID)
	assert.NoError(t, err)

	err = service.Activate(BudgetID{uuid.New()}, AccountID{uuid.New()})
	assert.EqualError(t, err, ErrAccountListDoesNotExist.Error())

	id, err := service.Add(budgetID, "new", MoneyAmount(1300))
	assert.NoError(t, err)
	err = service.Activate(budgetID, id)
	assert.NoError(t, err)
}

func TestAccountService_Cancel(t *testing.T) {
	budgetID := BudgetID{uuid.New()}
	repo := &mockRepo{make(map[BudgetID]*AccountList)}
	service := NewAccountListService(repo)
	err := service.Create(budgetID)
	assert.NoError(t, err)

	err = service.Cancel(BudgetID{uuid.New()}, AccountID{uuid.New()})
	assert.EqualError(t, err, ErrAccountListDoesNotExist.Error())

	id, err := service.Add(budgetID, "new", MoneyAmount(1300))
	assert.NoError(t, err)
	err = service.Cancel(budgetID, id)
	assert.NoError(t, err)
}

func TestAccountService_Delete(t *testing.T) {
	budgetID := BudgetID{uuid.New()}
	repo := &mockRepo{make(map[BudgetID]*AccountList)}
	service := NewAccountListService(repo)
	err := service.Create(budgetID)
	assert.NoError(t, err)

	err = service.Delete(BudgetID{uuid.New()}, AccountID{uuid.New()})
	assert.EqualError(t, err, ErrAccountListDoesNotExist.Error())

	id, err := service.Add(budgetID, "new", MoneyAmount(1300))
	assert.NoError(t, err)
	err = service.Delete(budgetID, id)
	assert.NoError(t, err)
}
