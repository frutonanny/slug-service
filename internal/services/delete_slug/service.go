package delete_slug

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/frutonanny/slug-service/internal/services"
)

type slugRepository interface {
	Delete(ctx context.Context, name string) (int64, error)
}

type userSlugsRepository interface {
	GetUserSlug(ctx context.Context, userID uuid.UUID) ([]string, error)
	DeleteUserSlugBySlugID(ctx context.Context, id int64) ([]uuid.UUID, error)
}

type operationRepository interface {
	AddOperation(ctx context.Context, userID uuid.UUID, slugID int64, event string) error
}

type transactor interface {
	RunInTx(ctx context.Context, f func(context.Context) error) error
}

type Service struct {
	log *zap.Logger

	slugRepository      slugRepository
	userSlugRepository  userSlugsRepository
	operationRepository operationRepository
	transactor          transactor
}

func New(
	log *zap.Logger,
	slugRepository slugRepository,
	userSlugRepository userSlugsRepository,
	operationRepository operationRepository,
	transactor transactor,
) *Service {
	return &Service{
		log:                 log,
		slugRepository:      slugRepository,
		userSlugRepository:  userSlugRepository,
		operationRepository: operationRepository,
		transactor:          transactor,
	}
}

// DeleteSlug - метод удаления, который делает следующие действия.
// 1. Помечает удаленным slug, чтобы все же хранился в истории в бд.
// 2. Удаляет у всех пользователей этот slug.
// 3. Заносит в историю операций запись об удалении slug у пользователей.
func (s *Service) DeleteSlug(ctx context.Context, name string) error {
	if err := s.transactor.RunInTx(ctx, func(ctx context.Context) error {
		id, err := s.slugRepository.Delete(ctx, name)
		if err != nil {
			s.log.Error("delete slug", zap.Error(err))
			return fmt.Errorf("delete slug: %v", err)
		}

		users, err := s.userSlugRepository.DeleteUserSlugBySlugID(ctx, id)
		if err != nil {
			s.log.Error("delete user slug by slugID", zap.Error(err))
			return fmt.Errorf("delete user slug by slugID: %v", err)
		}

		for _, userID := range users {
			if err := s.operationRepository.AddOperation(ctx, userID, id, services.DeleteSlug); err != nil {
				s.log.Error("add operation", zap.Error(err))
				return fmt.Errorf("add operation: %v", err)
			}
		}

		return nil
	}); err != nil {
		s.log.Error("run in tx", zap.Error(err))
		return fmt.Errorf("run in tx: %v", err)
	}

	return nil
}
