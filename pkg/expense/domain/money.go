package domain

import "errors"

var (
	InvalidCurrencyError = errors.New("currency is invalid")
)

type Currency string

type MoneyAmount struct {
	Amount   int
	Currency Currency
}

var availableCurrencies = map[Currency]bool{
	"RUB": true,
	"USD": true,
	"EUR": true,
}

func validateCurrency(c Currency) error {
	if _, ok := availableCurrencies[c]; !ok {
		return InvalidCurrencyError
	}
	return nil
}

func NewMoneyAmount(amount int, currency Currency) (MoneyAmount, error) {
	if err := validateCurrency(currency); err != nil {
		return MoneyAmount{}, err
	}
	return MoneyAmount{amount, currency}, nil
}
