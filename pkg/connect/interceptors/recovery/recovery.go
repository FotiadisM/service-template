package recovery

import (
	"context"
	"fmt"
	"runtime/debug"
)

type PanicError struct {
	Panic any
	Stack []byte
}

func (e *PanicError) Error() string {
	return fmt.Sprintf("panic caught: %v\n\n%s", e.Panic, e.Stack)
}

// DefaultRecoveryFunc will recover form panic p and return an err of type PanicError.
func DefaultRecoveryFunc(_ context.Context, p any) error {
	return &PanicError{Panic: p, Stack: debug.Stack()}
}
