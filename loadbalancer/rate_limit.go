package main

import (
	"net/http"
	"sync"
	"time"
)

type RateLimiter struct {
	visitors map[string]*Visitor
	mu       sync.Mutex
	rate     int           // tokens per interval
	interval time.Duration // refill interval
}

type Visitor struct {
	tokens int
	last   time.Time
}

func NewRateLimiter(rate int, interval time.Duration) *RateLimiter {
	return &RateLimiter{
		visitors: make(map[string]*Visitor),
		rate:     rate,
		interval: interval,
	}
}

func (rl *RateLimiter) getVisitor(ip string) *Visitor {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.visitors[ip]
	if !exists {
		v = &Visitor{tokens: rl.rate, last: time.Now()}
		rl.visitors[ip] = v
	}

	// refill tokens
	elapsed := time.Since(v.last)
	if elapsed > rl.interval {
		v.tokens = rl.rate
		v.last = time.Now()
	}

	return v
}

func (rl *RateLimiter) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		v := rl.getVisitor(ip)

		if v.tokens <= 0 {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}

		v.tokens--
		next.ServeHTTP(w, r)
	})
}
