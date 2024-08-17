package limiter

import (
	"context"
	"rate-limiter/validate"
	"sync"
	"time"
)

type (
	RateLimiter interface {
		Start(ctx context.Context)
		GetOutdatedKeys() []string
		Flush(keys []string)
		Limit(ip, handler string) bool
	}

	rateLimiter struct {
		frame     map[string]*Frame
		handlers  map[string]int
		whiteList map[string]struct{}

		gcInterval time.Duration
		lifetime   time.Duration

		mu sync.RWMutex
	}

	Frame struct {
		counter  int
		lifetime int64
	}
)

func NewRateLimiter(handlers map[string]int, whiteList map[string]struct{}, lifetime, gcInterval time.Duration) (RateLimiter, error) {
	var err error
	if err = validate.ValidateHandlers(handlers); err != nil {
		return nil, err
	}

	if err = validate.ValidateWhiteList(whiteList); err != nil {
		return nil, err
	}

	return &rateLimiter{
		frame:     make(map[string]*Frame),
		handlers:  handlers,
		whiteList: whiteList,

		gcInterval: gcInterval,
		lifetime:   lifetime,

		mu: sync.RWMutex{},
	}, nil
}

func (r *rateLimiter) Start(ctx context.Context) {
	NewGarbageCollector(r.gcInterval, r).Start(ctx)
}

func (r *rateLimiter) GetOutdatedKeys() []string {
	defer r.mu.Unlock()
	r.mu.Lock()

	keys := make([]string, 0)
	for k, v := range r.frame {
		if v.lifetime <= time.Now().Unix() {
			keys = append(keys, k)
		}
	}

	return keys
}

func (r *rateLimiter) Flush(keys []string) {
	defer r.mu.Unlock()
	r.mu.Lock()

	for _, v := range keys {
		delete(r.frame, v)
	}
}

func (r *rateLimiter) Limit(ip, handler string) bool {
	defer r.mu.Unlock()
	r.mu.Lock()

	var ok bool
	if _, ok = r.whiteList[ip]; ok {
		return true
	}

	var ttlLimit int
	if ttlLimit, ok = r.handlers[handler]; !ok {
		return false
	}

	if _, ok = r.frame[ip+handler]; !ok {
		r.frame[ip+handler] = &Frame{
			lifetime: time.Now().Add(r.lifetime).Unix(),
		}
		return true
	}

	if r.frame[ip+handler].counter > ttlLimit {
		return false
	}

	r.frame[ip+handler].counter += 1

	return true
}
