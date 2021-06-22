package domain

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewMoneyAmount(t *testing.T) {
	_, err := NewMoneyAmount(0, "non-existent")
	assert.True(t, errors.Is(err, ErrorCurrencyInvalid))

	amount, err := NewMoneyAmount(4200, "USD")
	assert.NoError(t, err)
	assert.Equal(t, 4200, amount.Amount)
	assert.Equal(t, "USD", string(amount.Currency))
}
