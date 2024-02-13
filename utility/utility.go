package utility

import "os"

// input: no parameter
// output: string
// func KeyForCurrencyExchangeOpenExchange will return access key for exchange rates api
func KeyForCurrencyExchangeOpenExchange() string {
	return os.Getenv("CURRENCY_EXCHANGE_KEY_OPENEXCHANGERATES")
}
