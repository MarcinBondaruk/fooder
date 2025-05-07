package login_limiter

import (
	"strconv"
	"testing"
	"time"
)

func TestRegisterBelowLimit(t *testing.T) {
	limiter := NewLoginLimiter(3, 2*time.Second)

	ip := "127.0.0.1"

	for i := 0; i < 3; i++ {
		allowed := limiter.Register(ip)
		if !allowed {
			t.Errorf("expected attempt %d to be allowed", i+1)
		}
	}

	// 4th should be blocked
	if limiter.Register(ip) {
		t.Error("expected 4th attempt to be blocked")
	}
}

func TestCleanupRemovesOldEntries(t *testing.T) {
	limiter := NewLoginLimiter(3, 500*time.Millisecond)

	ip := "10.0.0.1"
	_ = limiter.Register(ip)

	time.Sleep(600 * time.Millisecond)

	limiter.cleanup()

	limiter.mu.Lock()
	_, exists := limiter.attempts[ip]
	limiter.mu.Unlock()

	if exists {
		t.Error("expected old IP entry to be cleaned up")
	}
}

func TestCleanerRoutineStops(t *testing.T) {
	limiter := NewLoginLimiter(3, 300*time.Millisecond)

	ip := "192.168.0.1"
	_ = limiter.Register(ip)

	stopCh := make(chan struct{})
	limiter.StartCleaner(stopCh)

	time.Sleep(500 * time.Millisecond)
	close(stopCh)
 
	time.Sleep(100 * time.Millisecond)

	limiter.mu.Lock()
	_, exists := limiter.attempts[ip]
	limiter.mu.Unlock()

	if exists {
		t.Error("expected IP to be removed by cleaner")
	}
}

func TestRegisterAfterTimeout(t *testing.T) {
	limiter := NewLoginLimiter(1, 300*time.Millisecond)

	ip := "8.8.8.8"

	allowed := limiter.Register(ip)
	if !allowed {
		t.Error("first attempt should be allowed")
	}

	// Should now be blocked
	allowed = limiter.Register(ip)
	if allowed {
		t.Error("second attempt should be blocked")
	}

	time.Sleep(350 * time.Millisecond)
	limiter.cleanup()

	// Should be allowed again after timeout and cleanup
	allowed = limiter.Register(ip)
	if !allowed {
		t.Error("expected new attempt after timeout to be allowed")
	}
}

func BenchmarkRegisterSingleIP(b *testing.B) {
	limiter := NewLoginLimiter(1000, time.Minute)
	ip := "127.0.0.1"

	for i := 0; i < b.N; i++ {
		limiter.Register(ip)
	}
}

func BenchmarkRegisterManyIPs(b *testing.B) {
	limiter := NewLoginLimiter(5, time.Minute)

	for i := 0; i < b.N; i++ {
		ip := "192.168.1." + strconv.Itoa(i%1000)
		limiter.Register(ip)
	}
}

func BenchmarkCleanupSmall(b *testing.B) {
	limiter := NewLoginLimiter(5, time.Millisecond)

	// dodaj trochę danych
	for i := 0; i < 100; i++ {
		ip := "10.0.0." + strconv.Itoa(i)
		limiter.Register(ip)
	}

	time.Sleep(10 * time.Millisecond)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		limiter.cleanup()
	}
}

func BenchmarkCleanupLarge(b *testing.B) {
	limiter := NewLoginLimiter(5, time.Millisecond)

	// 10k IP
	for i := 0; i < 10_000; i++ {
		ip := "10.0.0." + strconv.Itoa(i)
		limiter.Register(ip)
	}

	time.Sleep(10 * time.Millisecond)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		limiter.cleanup()
	}
}

func BenchmarkRegisterParallel(b *testing.B) {
	limiter := NewLoginLimiter(100, time.Minute)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {

			ip := "10.0.0." + strconv.Itoa(int(time.Now().UnixNano()%10000))
			limiter.Register(ip)
		}
	})
}
