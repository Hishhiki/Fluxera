package middleware

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"
)

type RateLimiter interface {
	Increment(ctx context.Context, key string, ttl time.Duration) (int64, error)
}

func RateLimit(cache RateLimiter, prefix string, limit int64, window time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := r.RemoteAddr
			if host, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
				ip = host
			}

			key := fmt.Sprintf("ratelimit:%s:%s", prefix, ip)

			count, err := cache.Increment(r.Context(), key, window)
			if err != nil {
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}

			if count > limit {
				http.Error(w, "too many requests", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
