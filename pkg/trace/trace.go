package trace

import (
	"context"
	"time"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type Func func(context.Context, time.Duration, error)

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Do calls the trace function with context and time elapsed
func Do(ctx context.Context, fn Func, since time.Time) {
	elapsed := time.Since(since).Truncate(time.Millisecond)
	if fn != nil {
		fn(ctx, elapsed, nil)
	}
}

// Err calls the trace function with an error
func Err(ctx context.Context, fn Func, err error) {
	if fn != nil && err != nil {
		fn(ctx, 0, err)
	}
}
