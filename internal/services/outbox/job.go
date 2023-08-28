package outbox

import (
	"context"
	"time"
)

const (
	defaultExecutionTimeout = 30 * time.Second
)

type Job interface {
	Name() string

	Handle(ctx context.Context, data string) error

	// ExecutionTimeout — время, предоставленное обработчику очереди для выполнения задачи.
	ExecutionTimeout() time.Duration
}

// DefaultJob нужна для встраивания в другие jobs.
type DefaultJob struct{}

func (j DefaultJob) ExecutionTimeout() time.Duration {
	return defaultExecutionTimeout
}
