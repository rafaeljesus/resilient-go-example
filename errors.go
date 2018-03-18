package srv

import "errors"

var (
	// ErrUserAlreadyExists indicates user already exists.
	ErrUserAlreadyExists = errors.New("user already exists")
	// ErrUserNotFound indicates user was not found.
	ErrUserNotFound = errors.New("user not found")
	// ErrCircuitBreakerOpen indicates the circuit breaker is on open state.
	ErrCircuitBreakerOpen = errors.New("circuit breaker is in open state")
)
