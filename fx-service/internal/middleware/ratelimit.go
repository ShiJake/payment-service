package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/twitchtv/twirp"
)

// RateLimiter implements a token bucket rate limiter per IP address
type RateLimiter struct {
	mu       sync.Mutex
	buckets  map[string]*bucket
	rate     int           // requests per window
	window   time.Duration // time window
	cleanupInterval time.Duration
}

type bucket struct {
	tokens   int
	lastSeen time.Time
}

// NewRateLimiter creates a new rate limiter
// rate: number of requests allowed per window
// window: time window duration
func NewRateLimiter(rate int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		buckets:  make(map[string]*bucket),
		rate:     rate,
		window:   window,
		cleanupInterval: 5 * time.Minute,
	}
	
	// Start cleanup goroutine to remove stale buckets
	go rl.cleanup()
	
	return rl
}

// Allow checks if a request from the given IP is allowed
func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	b, exists := rl.buckets[ip]
	
	if !exists {
		// New IP, create bucket with full tokens
		rl.buckets[ip] = &bucket{
			tokens:   rl.rate - 1,
			lastSeen: now,
		}
		return true
	}

	// Calculate elapsed time and refill tokens
	elapsed := now.Sub(b.lastSeen)
	tokensToAdd := int(elapsed / rl.window * time.Duration(rl.rate))
	
	b.tokens += tokensToAdd
	if b.tokens > rl.rate {
		b.tokens = rl.rate
	}
	b.lastSeen = now

	// Check if we have tokens available
	if b.tokens > 0 {
		b.tokens--
		return true
	}

	return false
}

// cleanup removes stale buckets periodically
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(rl.cleanupInterval)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for ip, b := range rl.buckets {
			if now.Sub(b.lastSeen) > rl.cleanupInterval {
				delete(rl.buckets, ip)
			}
		}
		rl.mu.Unlock()
	}
}

// RateLimitMiddleware returns an HTTP middleware that enforces rate limiting
func RateLimitMiddleware(limiter *RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract IP address
			ip := r.RemoteAddr
			if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
				ip = forwarded
			}

			// Check rate limit
			if !limiter.Allow(ip) {
				// Return Twirp error for rate limit exceeded
				twirp.WriteError(w, twirp.NewError(twirp.ResourceExhausted, "rate limit exceeded"))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

