package domain

import (
	"errors"
	"fmt"
)

var (
	ErrorCurrencyInvalid = errors.New("currency is invalid")
)

type Currency string

type MoneyAmount struct {
	Amount   int
	Currency Currency
}

var availableCurrencies = map[Currency]struct{}{
	"RUB": {},
	"USD": {},
	"EUR": {},
}

func validateCurrency(c Currency) error {
	if _, ok := availableCurrencies[c]; !ok {
		return fmt.Errorf("%w: %v", ErrorCurrencyInvalid, c)
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
