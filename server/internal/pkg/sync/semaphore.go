package sync

import (
	"context"
)

type Semaphore struct {
	limit chan struct{}
}

func NewSemaphore(size int) Semaphore {
	return Semaphore{
		limit: make(chan struct{}, size),
	}
}

func (s *Semaphore) Acquire(ctx context.Context) error {
	select {
	case s.limit <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (s *Semaphore) Release() {
	<-s.limit
}
