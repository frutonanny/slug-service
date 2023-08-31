//go:generate mockgen -source=service.go -destination=mocks/service.gen.go

package modify_slug

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	slugrepo "github.com/frutonanny/slug-service/internal/repositories/slug"
	"github.com/frutonanny/slug-service/internal/services"
)

var ErrSlugNotFound = errors.New("slug not found")

type Slug struct {
	Name string
	Ttl  time.Time
}

type usersRepo interface {
	CreateUserIfNotExist(ctx context.Context, userID uuid.UUID) error
	AddUserSlug(ctx context.Context, user uuid.UUID, slugID int64, name string) (int64, error)
	AddUserSlugWithTtl(ctx context.Context, userID uuid.UUID, slugID int64, name string, ttl time.Time) (int64, error)
	DeleteUserSlug(ctx context.Context, user uuid.UUID, slugID int64) error
}

type slugRepo interface {
	GetID(ctx context.Context, name string) (int64, error)
}

type eventsRepo interface {
	AddEvent(ctx context.Context, userID uuid.UUID, slugID int64, event string) (int64, error)
}

type transactor interface {
	RunInTx(ctx context.Context, f func(context.Context) error) error
}

type Service struct {
	log *zap.Logger

	slugRepo   slugRepo
	usersRepo  usersRepo
	eventsRepo eventsRepo
	transactor transactor
}

func New(
	log *zap.Logger,
	slugRepo slugRepo,
	usersRepo usersRepo,
	eventsRepo eventsRepo,
	transactor transactor,
) *Service {
	return &Service{
		log:        log,
		slugRepo:   slugRepo,
		usersRepo:  usersRepo,
		eventsRepo: eventsRepo,
		transactor: transactor,
	}
}

// ModifySlugs добавляет/удаляет пользователя в/из сегментов.
func (s *Service) ModifySlugs(ctx context.Context, userID uuid.UUID, add []Slug, delete []string) error {
	if err := s.transactor.RunInTx(ctx, func(ctx context.Context) error {
		if err := s.usersRepo.CreateUserIfNotExist(ctx, userID); err != nil {
			s.log.Error("create user if not exist", zap.Error(err))
			return fmt.Errorf("create user if not exist: %v", err)
		}

		if add != nil {
			for _, slug := range add {
				slugID, err := s.slugRepo.GetID(ctx, slug.Name)
				if err != nil {
					if errors.Is(slugrepo.ErrRepoSlugNotFound, err) {
						return ErrSlugNotFound
					}

					s.log.Error("get id", zap.Error(err))
					return fmt.Errorf("get id: %v", err)
				}

				if _, err := s.eventsRepo.AddEvent(ctx, userID, slugID, services.AddSlug); err != nil {
					s.log.Error("add event", zap.Error(err))
					return fmt.Errorf("add event: %v", err)
				}

				if slug.Ttl.IsZero() {
					if _, err := s.usersRepo.AddUserSlug(ctx, userID, slugID, slug.Name); err != nil {
						s.log.Error("add users_slug", zap.Error(err))
						return fmt.Errorf("add users_slug: %v", err)
					}
				} else {
					if _, err := s.usersRepo.AddUserSlugWithTtl(ctx, userID, slugID, slug.Name, slug.Ttl); err != nil {
						s.log.Error("add users_slugs with ttl", zap.Error(err))
						return fmt.Errorf("add users_slugs with ttl: %v", err)
					}
				}
			}
		}

		if delete != nil {
			for _, name := range delete {
				slugID, err := s.slugRepo.GetID(ctx, name)
				if err != nil {
					if errors.Is(slugrepo.ErrRepoSlugNotFound, err) {
						return ErrSlugNotFound
					}

					s.log.Error("get id", zap.Error(err))
					return fmt.Errorf("get id: %v", err)
				}

				if err := s.usersRepo.DeleteUserSlug(ctx, userID, slugID); err != nil {
					s.log.Error("delete users_slugs", zap.Error(err))
					return fmt.Errorf("delete users_slugs: %v", err)
				}

				if _, err := s.eventsRepo.AddEvent(ctx, userID, slugID, services.DeleteSlug); err != nil {
					s.log.Error("add events", zap.Error(err))
					return fmt.Errorf("add event: %v", err)
				}
			}
		}

		return nil
	}); err != nil {
		s.log.Error("run in tx", zap.Error(err))
		return fmt.Errorf("run in tx: %v", err)
	}

	return nil
}
