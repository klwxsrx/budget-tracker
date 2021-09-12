package domain

import (
	"errors"
	"fmt"
)

var (
	ErrCurrencyInvalid = errors.New("currency is invalid")
)

type Currency string

type MoneyAmount struct {
	Amount   int
	Currency Currency
}

// nolint:gochecknoglobals
var availableCurrencies = map[Currency]struct{}{ // TODO: delete currency
	"RUB": {},
	"USD": {},
	"EUR": {},
}

func validateCurrency(c Currency) error {
	if _, ok := availableCurrencies[c]; !ok {
		return fmt.Errorf("%w: %v", ErrCurrencyInvalid, c)
	}
	return nil
}

func NewMoneyAmount(amount int, currency string) (MoneyAmount, error) {
	domainCurrency := Currency(currency)
	if err := validateCurrency(domainCurrency); err != nil {
		return MoneyAmount{}, err
	}
	return MoneyAmount{amount, domainCurrency}, nil
}
