package middleware

import (
	"math/rand/v2"
	"net/http"
	"time"

	"github.com/twitchtv/twirp"
)

const (
	maxJitterMs       = 500
	errorRate         = 10 // 1 in 10 requests
	slowRequestRate   = 10 // 1 in 10 requests
	slowRequestDelay  = 10 * time.Second
	slowRequestJitter = 2 * time.Second
)

// ChaosMiddleware simulates network jitter and periodic failures
func ChaosMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate network jitter
		jitter := time.Duration(rand.IntN(maxJitterMs)) * time.Millisecond
		time.Sleep(jitter)

		// Randomly fail requests
		if rand.IntN(errorRate) == 0 {
			err := twirp.WriteError(w, twirp.InternalError("service temporarily unavailable"))
			if err != nil {
				return
			}
			return
		}

		// Randomly slow down requests
		if rand.IntN(slowRequestRate) == 0 {
			delay := slowRequestDelay + time.Duration(rand.IntN(int(slowRequestJitter.Milliseconds())))*time.Millisecond
			time.Sleep(delay)
		}

		next.ServeHTTP(w, r)
	})
}
