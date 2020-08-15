package domain

type Currency string

var availableCurrencies = map[Currency]bool{
	"RUB": true,
	"USD": true,
	"EUR": true,
}

func IsCurrencyAvailable(currency Currency) bool {
	_, ok := availableCurrencies[currency]
	return ok
}
