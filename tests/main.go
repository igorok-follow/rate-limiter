package main

import (
	"context"
	"log"
	"math/rand"
	"rate-limiter/limiter"
	"time"
)

func main() {
	addresses := map[int]string{
		0: "1.1.1.1",
		1: "8.8.8.8",
		2: "8.8.4.4",
		3: "192.168.0.1",
	}

	handlers := map[string]int{
		"get_route":  5,
		"post_route": 5,
	}
	whiteList := map[string]struct{}{
		"192.168.0.1": {},
	}

	ctx, cancel := context.WithCancel(context.Background())
	r, err := limiter.NewRateLimiter(handlers, whiteList, time.Second, time.Second)
	if err != nil {
		log.Fatalln(err)
	}
	r.Start(ctx)

	for i := 0; i < 100; i++ {
		addr := addresses[rand.Intn(len(addresses))]
		log.Println(addr, r.Limit(addr, "get_route"))
		time.Sleep(time.Millisecond * 20)
	}

	cancel()
}
