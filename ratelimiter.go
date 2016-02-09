package ratelimiter

import (
	"time"
)

type RateLimiter struct {
	capacity int
	interval time.Duration

	ticker   *time.Ticker
	cancel   chan bool
	returned chan bool

	available chan bool
}

func New(capacity int, interval time.Duration) *RateLimiter {
	return &RateLimiter{
		capacity:  capacity,
		interval:  interval,
		cancel:    make(chan bool),
		returned:  make(chan bool, capacity),
		available: make(chan bool, capacity),
	}
}

func (d *RateLimiter) Start() {
	d.ticker = time.NewTicker(d.interval)

	for i := 0; i < d.capacity; i++ {
		d.available <- true
	}

	go func() {
		for {
			select {
			case <-d.ticker.C:
				d.refillAvailable()

			case <-d.cancel:
				d.ticker.Stop()
				return
			}
		}
	}()
}

func (d *RateLimiter) refillAvailable() {
	for {
		select {
		case <-d.returned:
			d.available <- true
		default:
			return
		}
	}
}

func (d *RateLimiter) Stop() {
	d.cancel <- true
}

func (d *RateLimiter) GetCapacity(n int) {
	for i := 0; i < n; i++ {
		<-d.available
	}
}

func (d *RateLimiter) ReleaseCapacity(n int) {
	for i := 0; i < n; i++ {
		d.returned <- true
	}
}
