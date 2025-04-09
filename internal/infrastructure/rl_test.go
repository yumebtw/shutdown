package infrastructure

import (
	"testing"
	"time"
)

func TestRateLimiter_Allow(t *testing.T) {
	start := time.Now()

	rps := 10
	rl := NewRateLimiter(rps)

	for i := 1; i <= rps; i++ {
		if !rl.Allow() {
			t.Error("allow failed for attempt", i, "time since start:", time.Since(start))
		}
	}

	if rl.Allow() {
		t.Error("allow should not have worked for attempt", rps+1)
	}

	time.Sleep(time.Second / time.Duration(rps))
	if !rl.Allow() {
		t.Error("allow failed for attempt", rps+2)
	}
}
