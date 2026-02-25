package nba

import (
	"sync"
	"time"
)

// RateLimiter enforces a minimum interval between API requests.
type RateLimiter struct {
	mu              sync.Mutex
	lastRequestTime time.Time
	minInterval     time.Duration
}

// NewRateLimiter creates a rate limiter with the given minimum interval.
func NewRateLimiter(minInterval time.Duration) *RateLimiter {
	if minInterval < 0 {
		minInterval = 0
	}
	return &RateLimiter{minInterval: minInterval}
}

// Wait blocks until the minimum interval has elapsed since the last request.
func (rl *RateLimiter) Wait() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	elapsed := time.Since(rl.lastRequestTime)
	if elapsed < rl.minInterval {
		time.Sleep(rl.minInterval - elapsed)
	}
	rl.lastRequestTime = time.Now()
}
