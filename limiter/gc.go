package limiter

import (
	"context"
	"log"
	"time"
)

type (
	Gc interface {
		Start(ctx context.Context)
	}

	gc struct {
		interval    time.Duration
		rateLimiter RateLimiter
	}
)

func NewGarbageCollector(interval time.Duration, limiter RateLimiter) Gc {
	return &gc{
		interval:    interval,
		rateLimiter: limiter,
	}
}

func (g *gc) Start(ctx context.Context) {
	go func() {
		for {
			select {
			case <-time.After(g.interval):
				keys := g.rateLimiter.GetOutdatedKeys()
				g.rateLimiter.Flush(keys)
				log.Println("flushed")
			case <-ctx.Done():
				return
			}
		}
	}()
}
