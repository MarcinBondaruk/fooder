package login_limiter

import (
	"sync"
	"time"
)

type LoginLimiter struct {
	mu       sync.Mutex
	limit    int
	timeout  time.Duration
	attempts map[string][]time.Time
}

func NewLoginLimiter(limit int, timeout time.Duration) *LoginLimiter {
	return &LoginLimiter{
		attempts: make(map[string][]time.Time),
		limit:    limit,
		timeout:  timeout,
	}
}

func (ll *LoginLimiter) Register(ip string) bool {
	ll.mu.Lock()
	defer ll.mu.Unlock()

	timestamp := time.Now()
	attempts := ll.attempts[ip]

	if len(attempts) >= ll.limit {
		return false
	}

	ll.attempts[ip] = append(attempts, timestamp)
	return true
}

func (ll *LoginLimiter) StartCleaner(stopCh <-chan struct{}) {
	ticker := time.NewTicker(ll.timeout)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				ll.cleanup()
			case <-stopCh:
				return
			}
		}
	}()
}

func (ll *LoginLimiter) cleanup() {
	ll.mu.Lock()
	defer ll.mu.Unlock()

	cutoff := time.Now().Add(-ll.timeout)

	for ip, times := range ll.attempts {
		if len(times) == 0 {
			delete(ll.attempts, ip)
			continue
		}

		if times[len(times)-1].Before(cutoff) {
			delete(ll.attempts, ip)
		}
	}
}
