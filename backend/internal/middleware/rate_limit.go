package middleware

import (
	"net/http"
	"sync"
	"time"
)

type RateLimiter struct {
	mu       sync.Mutex
	clients  map[string][]time.Time
	limit    int
	duration time.Duration
}

func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		clients:  make(map[string][]time.Time),
		limit:    100,
		duration: time.Minute,
	}
}

func (rl *RateLimiter) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
			ip = xff
		} else if xri := r.Header.Get("X-Real-IP"); xri != "" {
			ip = xri
		}

		rl.mu.Lock()

		now := time.Now()
		windowStart := now.Add(-rl.duration)

		var validRequests []time.Time
		for _, t := range rl.clients[ip] {
			if t.After(windowStart) {
				validRequests = append(validRequests, t)
			}
		}

		if len(validRequests) >= rl.limit {
			rl.clients[ip] = validRequests
			rl.mu.Unlock()
			http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		validRequests = append(validRequests, now)
		rl.clients[ip] = validRequests
		rl.mu.Unlock()

		next.ServeHTTP(w, r)
	})
}
