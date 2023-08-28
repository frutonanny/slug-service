package delete_slug

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/frutonanny/slug-service/internal/services"
)

type userSlugsRepository interface {
	GetUserSlug(ctx context.Context, userID uuid.UUID) ([]string, error)
	DeleteUserSlugByTtl(ctx context.Context, userID uuid.UUID) ([]int64, error)
}

type operationRepository interface {
	AddOperation(ctx context.Context, userID uuid.UUID, slugID int64, event string) error
}

type Service struct {
	log *zap.Logger

	userSlugRepository  userSlugsRepository
	operationRepository operationRepository
}

func New(
	log *zap.Logger,
	userSlugRepository userSlugsRepository,
	operationRepository operationRepository,
) *Service {
	return &Service{
		log:                 log,
		userSlugRepository:  userSlugRepository,
		operationRepository: operationRepository,
	}
}

// GetUserSlug отдает список slug пользователя.
// 1. Запускаем процесс по удалению отработанных slugs, по-установленному ttl
// 2. Делаем запрос на получение всех slugs пользователя.
// 3. Записываем в историю операций об удалении slug, у которых сработал ttl.
func (s *Service) GetUserSlug(ctx context.Context, userID uuid.UUID) ([]string, error) {
	slugIDs, err := s.userSlugRepository.DeleteUserSlugByTtl(ctx, userID)
	if err != nil {
		s.log.Error("delete slug", zap.Error(err))
		return nil, fmt.Errorf("delete slug: %v", err)
	}

	result, err := s.userSlugRepository.GetUserSlug(ctx, userID)
	if err != nil {
		s.log.Error("get user slug", zap.Error(err))
		return nil, fmt.Errorf("delete slug: %v", err)
	}

	for _, slugID := range slugIDs {
		if err := s.operationRepository.AddOperation(ctx, userID, slugID, services.DeleteSlug); err != nil {
			s.log.Error("add operation", zap.Error(err))
			return nil, fmt.Errorf("add operation: %v", err)
		}
	}

	return result, nil
}
