package domain

import "errors"

var InvalidCurrency = errors.New("currency is invalid")

type Currency string

var availableCurrencies = map[Currency]bool{
	"RUB": true,
	"USD": true,
	"EUR": true,
}

func validateCurrency(c Currency) error {
	if !isCurrencyAvailable(c) {
		return InvalidCurrency
	}
	return nil
}

func isCurrencyAvailable(currency Currency) bool {
	_, ok := availableCurrencies[currency]
	return ok
}
