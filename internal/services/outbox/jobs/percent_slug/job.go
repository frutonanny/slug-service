package percent_slug

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"

	usersrepo "github.com/frutonanny/slug-service/internal/repositories/users"
	"github.com/frutonanny/slug-service/internal/services"
	"github.com/frutonanny/slug-service/internal/services/outbox"
)

const Name = "percent_slug"

type sortingHat interface {
	Hit(userID uuid.UUID, percent int64) bool
}

type usersRepo interface {
	GetNextUser(ctx context.Context, previousID int64) (usersrepo.User, error)
	AddUserSlug(ctx context.Context, userID uuid.UUID, slugID int64, name string) (int64, error)
}

type eventsRepo interface {
	AddEvent(ctx context.Context, userID uuid.UUID, slugID int64, event string) (int64, error)
}

type transactor interface {
	RunInTx(ctx context.Context, f func(context.Context) error) error
}

type Job struct {
	outbox.DefaultJob
	log        *zap.Logger
	sortingHat sortingHat
	usersRepo  usersRepo
	eventsRepo eventsRepo
	transactor transactor
}

func New(
	log *zap.Logger,
	sortingHat sortingHat,
	usersRepo usersRepo,
	eventsRepo eventsRepo,
	transactor transactor,
) *Job {
	return &Job{
		log:        log,
		sortingHat: sortingHat,
		usersRepo:  usersRepo,
		eventsRepo: eventsRepo,
		transactor: transactor,
	}
}

func (j *Job) Name() string {
	return Name
}

// В допущениях озвучить сомнения, что задача долгая, надо прробежаться по всем пользователям,
// и если что-то пойдет не так, то надо либо повторять (это надо специально ещё логику писать),
// или оставлять как есть неконсистентное состояние (у кого-то будет сегмент, а кого-то нет, потому что до них
// не дошли).

// В допущениях написать, что новые создаваемые пользователи не будут добавляться в уже существующие
// процентные сегменты. Это просто не успела сделать. Так бы добавила в создание пользователя новую фоновую задачу
// в аутбокс.

func (j *Job) Handle(ctx context.Context, rawData string) error {
	j.log.Info("handle job", zap.String("name", Name))

	data := Data{}
	if err := json.Unmarshal([]byte(rawData), &data); err != nil {
		return fmt.Errorf("unmarshal data: %v", err)
	}

	previousUser := usersrepo.User{}

	for {
		user, err := j.usersRepo.GetNextUser(ctx, previousUser.ID)
		if err != nil {
			if errors.Is(err, usersrepo.ErrUserNotFound) {
				// Пользователи закончились – выходим.
				break
			}

			return fmt.Errorf("get first user: %v", err)
		}

		if j.sortingHat.Hit(user.UserID, int64(data.Percent)) {
			if err := j.transactor.RunInTx(ctx, func(ctx context.Context) error {
				if _, err := j.usersRepo.AddUserSlug(ctx, user.UserID, data.SlugID, data.SlugName); err != nil {
					return fmt.Errorf("add user slug: %v", err)
				}

				if _, err := j.eventsRepo.AddEvent(ctx, user.UserID, data.SlugID, services.AddSlug); err != nil {
					return fmt.Errorf("add event: %v", err)
				}

				return nil
			}); err != nil {
				return fmt.Errorf("run in tx: %v", err)
			}
		}

		previousUser = user
	}

	return nil
}
