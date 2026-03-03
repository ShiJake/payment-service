# FX Service

Foreign exchange rate service providing currency conversion quotes.
**Note:** This service simulates real-world conditions including network latency, intermittent failures, and rate limiting.

## Prerequisites

- [Docker](https://www.docker.com/get-started)

## Getting Started

```bash
npm start
```

Service runs on `http://localhost:4000`

## API

### Get Supported Currencies

**Request:**
```bash
curl -X POST http://localhost:4000/twirp/payments.v1.FXService/GetSupportedCurrencies \
  -H "Content-Type: application/json" \
  -d '{}'
```

**Response:**
```json
{
  "currencies": ["USD", "EUR", "GBP", "JPY", "CAD", "AUD", "CHF", "CNY", "INR", "MXN"]
}
```

### Get Quote

**Request:**
```bash
curl -X POST http://localhost:4000/twirp/payments.v1.FXService/GetQuote \
  -H "Content-Type: application/json" \
  -d '{
    "source_currency": "USD",
    "target_currency": "EUR"
  }'
```

**Response:**
```json
{
  "exchange_rate": 0.9215,
  "expiry_time": "2024-11-22T14:32:00Z"
}
```