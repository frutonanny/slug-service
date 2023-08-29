//go:generate mockgen -source=service.go -destination=mocks/service.gen.go

package get_user_slug

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	usersrepo "github.com/frutonanny/slug-service/internal/repositories/users"
	"github.com/frutonanny/slug-service/internal/services"
)

type Slug struct {
	Name      string
	Ttl       time.Time
	DeletedAt time.Time
}

type usersRepo interface {
	GetUserSlugs(ctx context.Context, userID uuid.UUID) ([]usersrepo.Slug, error)
	DeleteUserSlug(ctx context.Context, user uuid.UUID, slugID int64) error
}

type eventsRepo interface {
	AddEvents(ctx context.Context, userID uuid.UUID, slugID int64, event string) (int64, error)
}

type transactor interface {
	RunInTx(ctx context.Context, f func(context.Context) error) error
}

type Service struct {
	log *zap.Logger

	usersRepo  usersRepo
	eventsRepo eventsRepo
	transactor transactor
}

func New(
	log *zap.Logger,
	usersRepo usersRepo,
	eventsRepo eventsRepo,
	transactor transactor,
) *Service {
	return &Service{
		log:        log,
		usersRepo:  usersRepo,
		eventsRepo: eventsRepo,
		transactor: transactor,
	}
}

// GetUserSlug отдает список slug пользователя.
func (s *Service) GetUserSlug(ctx context.Context, userID uuid.UUID) ([]string, error) {
	// Возвращает все сегменты пользователя с ttl и deleted_at.
	slugs, err := s.usersRepo.GetUserSlugs(ctx, userID)
	if err != nil {
		s.log.Error("get user slugs", zap.Error(err))
		return nil, fmt.Errorf("get user slugs: %v", err)
	}

	if len(slugs) == 0 {
		return []string{}, nil
	}

	var (
		slugsForDelete []int64
		result         []string
	)

	for _, slug := range slugs {
		// Сегмент был удален, добавляем в список на удаление, в результат не попадает.
		// Или ttl уже прошел, то же самое.
		if !slug.DeletedAt.IsZero() || !slug.Ttl.IsZero() && slug.Ttl.Before(time.Now()) {
			slugsForDelete = append(slugsForDelete, slug.ID)
			continue
		}

		result = append(result, slug.Name)
	}

	// Неблокирующее удаление удаленных и истекших slugs у пользователя.
	if len(slugsForDelete) > 0 {
		go func() {
			if err := s.transactor.RunInTx(ctx, func(ctx context.Context) error {
				for _, slugID := range slugsForDelete {
					if err := s.usersRepo.DeleteUserSlug(ctx, userID, slugID); err != nil {
						s.log.Error("delete users_slugs", zap.Error(err))
						return fmt.Errorf("delete users_slugs: %v", err)
					}

					if _, err := s.eventsRepo.AddEvents(
						ctx,
						userID,
						slugID,
						services.DeleteSlug,
					); err != nil {
						s.log.Error("add event", zap.Error(err))
						return fmt.Errorf("add event: %v", err)
					}
				}
				return nil
			}); err != nil {
				s.log.Error("run in tx", zap.Error(err))
				return
			}
		}()
	}

	return result, nil
}
