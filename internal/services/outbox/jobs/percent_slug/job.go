package percent_slug

import (
	"context"
	"github.com/frutonanny/slug-service/internal/services/outbox"
	"go.uber.org/zap"
)

const Name = "percent_slug"

type Job struct {
	outbox.DefaultJob // todo хочу переложить в поле
	log               *zap.Logger
}

func New(log *zap.Logger) *Job {
	return &Job{
		log: log,
	}
}

func (j *Job) Name() string {
	return Name
}

func (j *Job) Handle(ctx context.Context, data string) error {
	j.log.Info("job handled", zap.String("name job: ", Name))
	// todo
	return nil
}
