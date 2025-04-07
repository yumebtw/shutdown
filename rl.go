package main

import "time"

type RateLimiter struct {
	rps   int
	token chan struct{}
}

func NewRateLimiter(rps int) *RateLimiter {
	return &RateLimiter{
		token: make(chan struct{}, rps),
		rps:   rps,
	}
}

func (r *RateLimiter) Run() {
	go func() {
		ticker := time.NewTicker(time.Second / time.Duration(r.rps))
		defer ticker.Stop()

		for range ticker.C {
			select {
			case r.token <- struct{}{}:
			default:
			}
		}
	}()
}

func (r *RateLimiter) Allow() bool {
	select {
	case <-r.token:
		return true
	default:
		return false
	}
}
