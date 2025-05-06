package login_limiter

import (
	"sync"
	"time"
)

const defaultCleanupInterval = 5 * time.Minute

type LoginLimiter struct {
	mu       *sync.Mutex
	limit    int
	window   time.Duration
	attempts map[string][]string
}

func NewLoginLimiter(limit int, window time.Duration) *LoginLimiter {
	return &LoginLimiter{
		attempts: make(map[string][]string),
		limit:    limit,
		window:   window,
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

	attempts = append(attempts, timestamp.String())
	return true
}

func (ll *LoginLimiter) StartCleaner(interval time.Duration, stopCh <-chan struct{}) {
	if interval <= 0 {
		interval = defaultCleanupInterval
	}
	ticker := time.NewTicker(interval)
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

}
