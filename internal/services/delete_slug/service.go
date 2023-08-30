//go:generate mockgen -source=service.go -destination=mocks/service.gen.go

package delete_slug

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/zap"

	slugrepo "github.com/frutonanny/slug-service/internal/repositories/slug"
)

var ErrSlugNotFound = errors.New("slug not found")

type slugRepo interface {
	Delete(ctx context.Context, name string) error
}

type Service struct {
	log *zap.Logger

	slugRepo slugRepo
}

func New(
	log *zap.Logger,
	slugRepo slugRepo,
) *Service {
	return &Service{
		log:      log,
		slugRepo: slugRepo,
	}
}

// DeleteSlug
// 1. Помечает сегмент удаленным.
// 2. Пишет событие в slug_events, что сегмент был удален.
func (s *Service) DeleteSlug(ctx context.Context, name string) error {
	if err := s.slugRepo.Delete(ctx, name); err != nil {
		if errors.Is(slugrepo.ErrRepoSlugNotFound, err) {
			return ErrSlugNotFound
		}

		s.log.Error("delete slug", zap.Error(err))
		return fmt.Errorf("delete slug: %v", err)
	}

	return nil
}
