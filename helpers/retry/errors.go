package retry

import (
	"fmt"
)

type ErrRetriesExceeded struct {
	tries int
	inner error
}

func NewErrRetriesExceeded(tries int, inner error) *ErrRetriesExceeded {
	return &ErrRetriesExceeded{
		tries: tries,
		inner: inner,
	}
}

func (e *ErrRetriesExceeded) Error() string {
	return fmt.Sprintf("limit of %d retries exceeded: %v", e.tries, e.inner)
}

func (e *ErrRetriesExceeded) Unwrap() error {
	return e.inner
}

func (e *ErrRetriesExceeded) Is(err error) bool {
	ee, ok := err.(*ErrRetriesExceeded)
	if !ok {
		return false
	}

	return ee.tries == e.tries && ee.inner == e.inner
}

func (e *ErrRetriesExceeded) Tries() int {
	return e.tries
}
