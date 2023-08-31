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
	AddEventWithCreatedAt(
		ctx context.Context,
		userID uuid.UUID,
		slugID int64,
		event string,
		createdAt time.Time,
	) (int64, error)
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

// GetUserSlugs отдает список сегментов пользователя.
func (s *Service) GetUserSlugs(ctx context.Context, userID uuid.UUID, sync bool) ([]string, error) {
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
		slugsForDelete []slugForDelete
		result         []string
	)

	// Проходим по всем сегментам и удаляем невалидные:
	for _, slug := range slugs {
		// Если сегмент был удален, то добавляем его в список на удаление (в результат не попадает).
		// Также если ttl (автоматическое удаление) уже прошло, то также добавляем в список на удаление.
		if !slug.DeletedAt.IsZero() {
			slugsForDelete = append(slugsForDelete, slugForDelete{
				id:        slug.ID,
				deletedAt: slug.DeletedAt,
			})
			continue
		}

		if !slug.Ttl.IsZero() && slug.Ttl.Before(time.Now()) {
			slugsForDelete = append(slugsForDelete, slugForDelete{
				id:        slug.ID,
				deletedAt: slug.Ttl,
			})
			continue
		}

		result = append(result, slug.Name)
	}

	// Блокирующее или неблокирующее удаление удаленных и истекших slugs у пользователя.
	if len(slugsForDelete) > 0 {
		if sync {
			if err := s.deleteSlugs(ctx, userID, slugsForDelete); err != nil {
				s.log.Error("sync delete slugs", zap.Error(err))
				return nil, fmt.Errorf("delete slugs: %v", err)
			}
		} else {
			go func() {
				if err := s.deleteSlugs(context.Background(), userID, slugsForDelete); err != nil {
					s.log.Error("async delete slugs", zap.Error(err))
					return
				}
			}()
		}
	}

	if len(result) == 0 {
		return []string{}, nil
	}

	return result, nil
}

type slugForDelete struct {
	id        int64
	deletedAt time.Time
}

func (s *Service) deleteSlugs(ctx context.Context, userID uuid.UUID, slugsForDelete []slugForDelete) error {
	if err := s.transactor.RunInTx(ctx, func(ctx context.Context) error {
		for _, slug := range slugsForDelete {
			if err := s.usersRepo.DeleteUserSlug(ctx, userID, slug.id); err != nil {
				s.log.Error("delete users_slugs", zap.Error(err))
				return fmt.Errorf("delete users_slugs: %v", err)
			}

			if _, err := s.eventsRepo.AddEventWithCreatedAt(
				ctx,
				userID,
				slug.id,
				services.DeleteSlug,
				slug.deletedAt,
			); err != nil {
				s.log.Error("add event", zap.Error(err))
				return fmt.Errorf("add event: %v", err)
			}
		}

		return nil
	}); err != nil {
		s.log.Error("run in tx", zap.Error(err))
		return fmt.Errorf("run in tx: %v", err)
	}

	return nil
}
