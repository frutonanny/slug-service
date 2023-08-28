package modify_slug

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type AddSlug struct {
	Name string
	Ttl  time.Time
}

type userSlugsRepository interface {
	AddUserSlug(ctx context.Context, user uuid.UUID, id int64, name string) error
	DeleteUserSlug(ctx context.Context, user uuid.UUID, ids int64, name string) error
}

type slugRepository interface {
	GetID(ctx context.Context, name string) (int64, error)
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
	userSlugsRepository userSlugsRepository
	operationRepository operationRepository
	transactor          transactor
}

func New(
	log *zap.Logger,
	slugRepository slugRepository,
	userSlugsRepository userSlugsRepository,
	operationRepository operationRepository,
	transactor transactor,
) *Service {
	return &Service{
		log:                 log,
		slugRepository:      slugRepository,
		userSlugsRepository: userSlugsRepository,
		operationRepository: operationRepository,
		transactor:          transactor,
	}
}

// ModifySlugs добавляет/удаляет пользователя в/из сегментов.
// Сценарий действий:
// Действия в транзакции:
//
//	Если нужно записать (add не пустой список):
//
// 1. Запрашиваем id-slugs.
//   - если есть незнакомые имена slugs - откатываем транзакцию и возвращаем ErrSlugNotFound;
//   - если имена slugs валидны, получаем id всех slugs.
//
// 2. Добавляем пользователю инфо о slugs:
//   - проверяем, передан ли ttl:
//   - если передан, то заносим инфу об отложенном удалении slug у пользователя;
//   - если не передан - просто записываем slug пользователю.
//     3. Записываем инфу о добавлении всех slugs в таблицу истории операций.
//     Если нужно удалить (delete не пустой список):
//   - запрашиваем id-slugs.
//   - если есть незнакомые имена slugs - откатываем транзакцию и возвращаем ErrSlugNotFound;
//   - если имена slugs валидны, получаем id всех slugs.
//     2. Удаляем у пользователя инфо о slugs:
//     3. Записываем инфу об удалении всех slugs в таблицу истории операций.
func (s *Service) ModifySlugs(ctx context.Context, userID uuid.UUID, add []AddSlug, delete []string) error {
	if err := s.transactor.RunInTx(ctx, func(ctx context.Context) error {
		if add != nil {

		}
		// TODO
		return nil
	}); err != nil {
		s.log.Error("run in tx", zap.Error(err))
		return fmt.Errorf("run in tx: %v", err)
	}

	return nil
}
