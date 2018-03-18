package main

import (
	"time"

	"github.com/rafaeljesus/resilient-go-example/http"
	"github.com/sony/gobreaker"
)

func main() {
	breaker := newCircuitBreaker()
	transport := http.NewTransport(breaker)
	req := http.NewRequest(
		http.WithRoundTripper(transport),
	)
	client := http.NewClient(req)

	_ = http.NewStoreService(client)

	// pass storer downstream to
	// http handler, messaging handler, etc...
}

func newCircuitBreaker() *gobreaker.CircuitBreaker {
	return gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:    "HTTP Client",
		Timeout: time.Second * 45,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= 3 && failureRatio >= 0.6
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			// do smth when circuit breaker trips.
		},
	})
}
