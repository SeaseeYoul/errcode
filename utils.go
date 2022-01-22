package errcode

import (
	"context"
	"time"
)

// Context get current context or a new context
func Context(e Codes) context.Context {
	return safeCode(e).Context()
}

// WithContext create a new context
func WithContext(e Codes, ctx context.Context) Codes {
	return safeCode(e).WithContext(ctx)
}

func WithCancel(e Codes) (codes Codes, cancel context.CancelFunc) {
	return safeCode(e).WithCancel()
}

func WithDeadline(e Codes, d time.Time) (Codes, context.CancelFunc) {
	return safeCode(e).WithDeadline(d)
}

func WithTimeout(e Codes, timeout time.Duration) (Codes, context.CancelFunc) {
	return safeCode(e).WithTimeout(timeout)
}

func WithValue(e Codes, key, val interface{}) Codes {
	return safeCode(e).WithValue(key, val)
}
