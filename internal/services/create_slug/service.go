//go:generate mockgen -source=service.go -destination=mocks/service.gen.go

package create_slug

import (
	"context"
	"encoding/json"
	"fmt"

	"go.uber.org/zap"

	slugrepo "github.com/frutonanny/slug-service/internal/repositories/slug"
	percentslugjob "github.com/frutonanny/slug-service/internal/services/outbox/jobs/percent_slug"
)

type Options struct {
	Percent *int
}

func (o *Options) IsEmpty() bool {
	return o.Percent == nil
}

type slugRepo interface {
	Create(ctx context.Context, name string, options slugrepo.Options) error
}

type outboxService interface {
	Put(ctx context.Context, name, data string) error
}

type transactor interface {
	RunInTx(ctx context.Context, f func(context.Context) error) error
}

type Service struct {
	log *zap.Logger

	outboxService outboxService
	slugRepo      slugRepo
	transactor    transactor
}

func New(
	log *zap.Logger,
	outboxService outboxService,
	slugRepo slugRepo,
	transactor transactor,
) *Service {
	return &Service{
		log:           log,
		outboxService: outboxService,
		slugRepo:      slugRepo,
		transactor:    transactor,
	}
}

func (s *Service) CreateSlug(ctx context.Context, name string, options Options) error {
	if !options.IsEmpty() {
		return s.createSlugWithOptions(ctx, name, options)
	}

	if err := s.slugRepo.Create(ctx, name, slugrepo.Options{}); err != nil {
		return fmt.Errorf("create slug: %v", err)
	}

	return nil
}

func (s *Service) createSlugWithOptions(ctx context.Context, name string, options Options) error {
	if err := s.transactor.RunInTx(ctx, func(ctx context.Context) error {
		if err := s.slugRepo.Create(
			ctx,
			name,
			slugrepo.Options{
				Percent: options.Percent,
			},
		); err != nil {
			s.log.Error("create slug", zap.Error(err))
			return fmt.Errorf("create slug: %v", err)
		}

		// Заготовка на случай, если появятся другие опции, требующие других фоновых задач.
		if options.Percent != nil {
			data := percentslugjob.Data{
				Name:    name,
				Percent: *options.Percent,
			}

			b, err := json.Marshal(data)
			if err != nil {
				s.log.Error("marshal percent slug job data", zap.Error(err))
				return fmt.Errorf("marshal percent slug job data: %v", err)
			}

			if err := s.outboxService.Put(ctx, percentslugjob.Name, string(b)); err != nil {
				s.log.Error("put to outbox", zap.Error(err))
				return fmt.Errorf("put to outbox: %v", err)
			}
		}

		return nil
	}); err != nil {
		s.log.Error("run in tx", zap.Error(err))
		return fmt.Errorf("run in tx: %v", err)
	}

	return nil
}
