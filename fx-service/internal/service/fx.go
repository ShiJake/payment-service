package service

import (
	"fmt"
	"math/rand"
	"time"

	"fxservice/internal/constants"
	pb "fxservice/rpc/fxservice"
)

// Quote represents a currency exchange quote
type Quote struct {
	ExchangeRate float64
	ExpiryTime   time.Time
}

// FXService handles foreign exchange operations
type FXService struct {
	rand *rand.Rand
}

// NewFXService creates a new FX service instance
func NewFXService() *FXService {
	return &FXService{
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// GetQuote calculates an exchange rate quote
func (s *FXService) GetQuote(sourceCurrency, targetCurrency pb.Currency) (*Quote, error) {
	// Validate currencies
	if !constants.IsCurrencySupported(sourceCurrency) {
		return nil, fmt.Errorf("unsupported source currency: %s", sourceCurrency)
	}

	if !constants.IsCurrencySupported(targetCurrency) {
		return nil, fmt.Errorf("unsupported target currency: %s", targetCurrency)
	}

	// Validate currency pair
	if !constants.IsPairSupported(sourceCurrency, targetCurrency) {
		return nil, fmt.Errorf("currency pair %s -> %s is not supported", sourceCurrency, targetCurrency)
	}

	// Calculate exchange rate with variation (±2%)
	sourceRate := constants.BaseRates[sourceCurrency]
	targetRate := constants.BaseRates[targetCurrency]
	baseExchangeRate := targetRate / sourceRate
	variation := 1.0 + (s.rand.Float64()*0.04 - 0.02)
	exchangeRate := baseExchangeRate * variation

	return &Quote{
		ExchangeRate: exchangeRate,
		ExpiryTime:   time.Now().Add(2 * time.Minute),
	}, nil
}

// GetSupportedCurrencies returns the list of supported currencies
func (s *FXService) GetSupportedCurrencies() []pb.Currency {
	return constants.SupportedCurrencies
}

