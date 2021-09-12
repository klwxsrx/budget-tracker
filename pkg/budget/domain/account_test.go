package domain

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"strings"
	"testing"
)

func TestAccountList_Add(t *testing.T) {
	listID := BudgetID{uuid.New()}
	list := LoadAccountList(&AccountListState{listID, nil})

	_, err := list.Add("", MoneyAmount{0, "USD"})
	assert.EqualError(t, err, ErrAccountInvalidTitle.Error())

	_, err = list.Add(" \n\t", MoneyAmount{0, "USD"})
	assert.EqualError(t, err, ErrAccountInvalidTitle.Error())

	_, err = list.Add(strings.Repeat("s", 101), MoneyAmount{0, "USD"})
	assert.EqualError(t, err, ErrAccountInvalidTitle.Error())

	id, err := list.Add("some", MoneyAmount{-42, "USD"})
	assert.NoError(t, err)
	assert.Len(t, list.GetChanges(), 1)
	createdEvent, ok := list.GetChanges()[0].(*AccountCreatedEvent)
	require.True(t, ok)
	assert.Equal(t, listID.String(), createdEvent.AggregateID().String())
	assert.Equal(t, accountListAggregateName, createdEvent.AggregateName())
	assert.Equal(t, EventTypeAccountCreated, createdEvent.Type())
	assert.Equal(t, id.String(), createdEvent.AccountID.String())
	assert.Equal(t, "some", createdEvent.Title)
	assert.Equal(t, -42, createdEvent.InitialBalance)
	assert.Equal(t, "USD", string(createdEvent.Currency))
	assert.Len(t, list.getAccounts(), 1)
	assert.Equal(t, id.String(), list.getAccounts()[0].GetID().String())
	assert.Equal(t, AccountStatusActive, list.getAccounts()[0].GetStatus())
	assert.Equal(t, "some", list.getAccounts()[0].GetTitle())
	assert.Equal(t, MoneyAmount{-42, "USD"}, list.getAccounts()[0].GetInitialBalance())

	_, err = list.Add("some", MoneyAmount{0, "USD"})
	assert.EqualError(t, err, ErrAccountDuplicateTitle.Error())

	_, err = list.Add("another", MoneyAmount{13, "RUB"})
	assert.NoError(t, err)
	assert.Len(t, list.GetChanges(), 2)
	assert.Len(t, list.getAccounts(), 2)

	_, err = list.Add(strings.Repeat("s", 100), MoneyAmount{0, "USD"})
	assert.NoError(t, err)
	assert.Len(t, list.GetChanges(), 3)
	assert.Len(t, list.getAccounts(), 3)
}

func TestAccountList_Reorder(t *testing.T) {
	assertOrder := func(list *AccountList, titles []string) {
		for i, acc := range list.getAccounts() {
			assert.Equal(t, titles[i], acc.GetTitle())
		}
	}
	listID := BudgetID{uuid.New()}
	id1 := AccountID{uuid.New()}
	id2 := AccountID{uuid.New()}
	id3 := AccountID{uuid.New()}
	id4 := AccountID{uuid.New()}
	id5 := AccountID{uuid.New()}
	list := LoadAccountList(&AccountListState{listID, []*AccountState{
		{id1, AccountStatusActive, "1", MoneyAmount{0, "USD"}},
		{id2, AccountStatusActive, "2", MoneyAmount{0, "USD"}},
		{id3, AccountStatusActive, "3", MoneyAmount{0, "USD"}},
		{id4, AccountStatusActive, "4", MoneyAmount{0, "USD"}},
		{id5, AccountStatusActive, "5", MoneyAmount{0, "USD"}},
	}})

	err := list.Reorder(AccountID{uuid.New()}, 3)
	assert.EqualError(t, err, ErrAccountDoesNotExist.Error())

	err = list.Reorder(id1, -10)
	assert.NoError(t, err)
	assert.Len(t, list.GetChanges(), 0)
	assertOrder(list, []string{"1", "2", "3", "4", "5"})

	err = list.Reorder(id1, 0)
	assert.NoError(t, err)
	assert.Len(t, list.GetChanges(), 0)
	assertOrder(list, []string{"1", "2", "3", "4", "5"})

	err = list.Reorder(id3, 2)
	assert.NoError(t, err)
	assert.Len(t, list.GetChanges(), 0)
	assertOrder(list, []string{"1", "2", "3", "4", "5"})

	err = list.Reorder(id5, 4)
	assert.NoError(t, err)
	assert.Len(t, list.GetChanges(), 0)
	assertOrder(list, []string{"1", "2", "3", "4", "5"})

	err = list.Reorder(id5, 10)
	assert.NoError(t, err)
	assert.Len(t, list.GetChanges(), 0)
	assertOrder(list, []string{"1", "2", "3", "4", "5"})

	err = list.Reorder(id1, 3)
	assert.NoError(t, err)
	assertOrder(list, []string{"2", "3", "4", "1", "5"})
	assert.Len(t, list.GetChanges(), 1)
	reorderedEvent, ok := list.GetChanges()[0].(*AccountReorderedEvent)
	require.True(t, ok)
	assert.Equal(t, listID.String(), reorderedEvent.AggregateID().String())
	assert.Equal(t, accountListAggregateName, reorderedEvent.AggregateName())
	assert.Equal(t, EventTypeAccountReordered, reorderedEvent.Type())
	assert.Equal(t, id1.String(), reorderedEvent.AccountID.String())
	assert.Equal(t, 3, reorderedEvent.Position)

	err = list.Reorder(id5, 1)
	assert.NoError(t, err)
	assertOrder(list, []string{"2", "5", "3", "4", "1"})
	assert.Len(t, list.GetChanges(), 2)

	err = list.Reorder(id3, 4)
	assert.NoError(t, err)
	assertOrder(list, []string{"2", "5", "4", "1", "3"})
	assert.Len(t, list.GetChanges(), 3)

	err = list.Reorder(id4, 0)
	assert.NoError(t, err)
	assertOrder(list, []string{"4", "2", "5", "1", "3"})
	assert.Len(t, list.GetChanges(), 4)

	err = list.Reorder(id2, 3)
	assert.NoError(t, err)
	assertOrder(list, []string{"4", "5", "1", "2", "3"})
	assert.Len(t, list.GetChanges(), 5)

	err = list.Reorder(id1, 1)
	assert.NoError(t, err)
	assertOrder(list, []string{"4", "1", "5", "2", "3"})
	assert.Len(t, list.GetChanges(), 6)

	err = list.Reorder(id5, 10)
	assert.NoError(t, err)
	assertOrder(list, []string{"4", "1", "2", "3", "5"})
	assert.Len(t, list.GetChanges(), 7)

	err = list.Reorder(id1, -10)
	assert.NoError(t, err)
	assertOrder(list, []string{"1", "4", "2", "3", "5"})
	assert.Len(t, list.GetChanges(), 8)
}

func TestAccountList_Rename(t *testing.T) {
	listID := BudgetID{uuid.New()}
	itemID := AccountID{uuid.New()}
	list := LoadAccountList(&AccountListState{listID, []*AccountState{
		{itemID, AccountStatusActive, "some", MoneyAmount{0, "USD"}},
		{AccountID{uuid.New()}, AccountStatusActive, "another", MoneyAmount{0, "USD"}},
	}})

	err := list.Rename(AccountID{uuid.New()}, "test")
	assert.EqualError(t, err, ErrAccountDoesNotExist.Error())

	err = list.Rename(itemID, " \t\n")
	assert.EqualError(t, err, ErrAccountInvalidTitle.Error())

	_, err = list.Add(strings.Repeat("s", 101), MoneyAmount{0, "USD"})
	assert.EqualError(t, err, ErrAccountInvalidTitle.Error())

	err = list.Rename(itemID, "another")
	assert.EqualError(t, err, ErrAccountDuplicateTitle.Error())

	err = list.Rename(itemID, "some")
	assert.NoError(t, err)
	assert.Len(t, list.GetChanges(), 0)

	err = list.Rename(itemID, "new")
	assert.NoError(t, err)
	assert.Len(t, list.GetChanges(), 1)
	assert.Len(t, list.getAccounts(), 2)
	assert.Equal(t, "new", list.getAccounts()[0].GetTitle())
	renamedEvent, ok := list.GetChanges()[0].(*AccountRenamedEvent)
	require.True(t, ok)
	assert.Equal(t, listID.String(), renamedEvent.AggregateID().String())
	assert.Equal(t, accountListAggregateName, renamedEvent.AggregateName())
	assert.Equal(t, EventTypeAccountRenamed, renamedEvent.Type())
	assert.Equal(t, itemID.String(), renamedEvent.AccountID.String())
	assert.Equal(t, "new", renamedEvent.Title)

	_, err = list.Add(strings.Repeat("s", 100), MoneyAmount{0, "USD"})
	assert.NoError(t, err)
}

func TestAccountList_Activate(t *testing.T) {
	listID := BudgetID{uuid.New()}
	item1ID := AccountID{uuid.New()}
	item2ID := AccountID{uuid.New()}
	list := LoadAccountList(&AccountListState{listID, []*AccountState{
		{item1ID, AccountStatusActive, "some", MoneyAmount{0, "USD"}},
		{item2ID, AccountStatusCancelled, "another", MoneyAmount{0, "USD"}},
	}})

	err := list.Activate(AccountID{uuid.New()})
	assert.EqualError(t, err, ErrAccountDoesNotExist.Error())

	err = list.Activate(item1ID)
	assert.NoError(t, err)
	assert.Len(t, list.GetChanges(), 0)

	err = list.Activate(item2ID)
	assert.NoError(t, err)
	assert.Len(t, list.GetChanges(), 1)
	event, ok := list.GetChanges()[0].(*AccountActivatedEvent)
	require.True(t, ok)
	assert.Equal(t, listID.String(), event.AggregateID().String())
	assert.Equal(t, accountListAggregateName, event.AggregateName())
	assert.Equal(t, EventTypeAccountActivated, event.Type())
	assert.Equal(t, item2ID.String(), event.AccountID.String())

	item2 := list.findAccount(item2ID)
	assert.NotNil(t, item2)
	assert.Equal(t, AccountStatusActive, item2.GetStatus())
}

func TestAccountList_Cancel(t *testing.T) {
	listID := BudgetID{uuid.New()}
	item1ID := AccountID{uuid.New()}
	item2ID := AccountID{uuid.New()}
	list := LoadAccountList(&AccountListState{listID, []*AccountState{
		{item1ID, AccountStatusCancelled, "some", MoneyAmount{0, "USD"}},
		{item2ID, AccountStatusActive, "another", MoneyAmount{0, "USD"}},
		{item2ID, AccountStatusActive, "third", MoneyAmount{0, "USD"}},
	}})

	err := list.Cancel(AccountID{uuid.New()})
	assert.EqualError(t, err, ErrAccountDoesNotExist.Error())

	err = list.Cancel(item1ID)
	assert.NoError(t, err)
	assert.Len(t, list.GetChanges(), 0)

	err = list.Cancel(item2ID)
	assert.NoError(t, err)
	assert.Len(t, list.GetChanges(), 1)
	event, ok := list.GetChanges()[0].(*AccountCancelledEvent)
	require.True(t, ok)
	assert.Equal(t, listID.String(), event.AggregateID().String())
	assert.Equal(t, accountListAggregateName, event.AggregateName())
	assert.Equal(t, EventTypeAccountCancelled, event.Type())
	assert.Equal(t, item2ID.String(), event.AccountID.String())

	item2 := list.findAccount(item2ID)
	assert.NotNil(t, item2)
	assert.Equal(t, AccountStatusCancelled, item2.GetStatus())
}

func TestAccountList_Delete(t *testing.T) {
	listID := BudgetID{uuid.New()}
	item1ID := AccountID{uuid.New()}
	list := LoadAccountList(&AccountListState{listID, []*AccountState{
		{item1ID, AccountStatusActive, "some", MoneyAmount{0, "USD"}},
	}})

	err := list.Delete(AccountID{uuid.New()})
	assert.EqualError(t, err, ErrAccountDoesNotExist.Error())

	err = list.Delete(item1ID)
	assert.NoError(t, err)
	assert.Len(t, list.getAccounts(), 0)
	assert.Len(t, list.GetChanges(), 1)
	event, ok := list.GetChanges()[0].(*AccountDeletedEvent)
	require.True(t, ok)
	assert.Equal(t, listID.String(), event.AggregateID().String())
	assert.Equal(t, accountListAggregateName, event.AggregateName())
	assert.Equal(t, EventTypeAccountDeleted, event.Type())
	assert.Equal(t, item1ID.String(), event.AccountID.String())
}
