package constants

import (
	pb "fxservice/rpc/fxservice"
)

// BaseRates contains exchange rates relative to USD
var BaseRates = map[pb.Currency]float64{
	pb.Currency_USD: 1.0,
	pb.Currency_EUR: 0.92,
	pb.Currency_GBP: 0.79,
	pb.Currency_JPY: 149.50,
	pb.Currency_CAD: 1.36,
	pb.Currency_AUD: 1.53,
	pb.Currency_CHF: 0.88,
	pb.Currency_CNY: 7.24,
	pb.Currency_INR: 83.12,
	pb.Currency_MXN: 17.15,
}

// SupportedCurrencies lists all valid currencies (excluding UNSPECIFIED)
var SupportedCurrencies = []pb.Currency{
	pb.Currency_USD,
	pb.Currency_EUR,
	pb.Currency_GBP,
	pb.Currency_JPY,
	pb.Currency_CAD,
	pb.Currency_AUD,
	pb.Currency_CHF,
	pb.Currency_CNY,
	pb.Currency_INR,
	pb.Currency_MXN,
}

// IsCurrencySupported checks if a currency is valid and supported
func IsCurrencySupported(currency pb.Currency) bool {
	if currency == pb.Currency_CURRENCY_UNSPECIFIED {
		return false
	}
	_, exists := BaseRates[currency]
	return exists
}

// IsPairSupported checks if a currency pair conversion is supported
// All non-UNSPECIFIED currency pairs are supported
func IsPairSupported(source, target pb.Currency) bool {
	return IsCurrencySupported(source) && IsCurrencySupported(target) && source != target
}

