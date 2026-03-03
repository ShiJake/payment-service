package handler

import (
	"context"
	"log"
	"net/http"

	"fxservice/internal/service"
	pb "fxservice/rpc/fxservice"

	"github.com/twitchtv/twirp"
)

// TwirpHandler implements the Twirp FXService interface
type TwirpHandler struct {
	svc *service.FXService
}

// NewTwirpHandler creates a new Twirp handler
func NewTwirpHandler(svc *service.FXService) pb.FXService {
	return &TwirpHandler{
		svc: svc,
	}
}

// GetQuote handles the GetQuote RPC call
func (h *TwirpHandler) GetQuote(ctx context.Context, req *pb.GetQuoteRequest) (*pb.GetQuoteResponse, error) {
	quote, err := h.svc.GetQuote(req.SourceCurrency, req.TargetCurrency)
	if err != nil {
		return nil, twirp.InvalidArgumentError("request", err.Error())
	}

	return &pb.GetQuoteResponse{
		ExchangeRate: quote.ExchangeRate,
		ExpiryTime:   quote.ExpiryTime.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

// GetSupportedCurrencies handles the GetSupportedCurrencies RPC call
func (h *TwirpHandler) GetSupportedCurrencies(ctx context.Context, req *pb.GetSupportedCurrenciesRequest) (*pb.GetSupportedCurrenciesResponse, error) {
	return &pb.GetSupportedCurrenciesResponse{
		Currencies: h.svc.GetSupportedCurrencies(),
	}, nil
}

// LoggingMiddleware logs all requests
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

