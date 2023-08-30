package modify_slug_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/frutonanny/slug-service/internal/services"
	modifyslug "github.com/frutonanny/slug-service/internal/services/modify_slug"
	mock_modify_slug "github.com/frutonanny/slug-service/internal/services/modify_slug/mocks"
)

func TestService_ModifySlugs(t *testing.T) {
	t.Run("success add slugs", func(t *testing.T) {
		// Arrange.
		ctx := context.Background()
		userID := uuid.New()
		name := "123"
		ttl := time.Now()
		adds := []modifyslug.Slug{
			{
				Name: name,
				Ttl:  ttl,
			},
		}

		slugID := int64(1)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		usersRepo := mock_modify_slug.NewMockusersRepo(ctrl)
		usersRepo.EXPECT().CreateUserIfNotExist(gomock.Any(), userID).Return(nil)
		usersRepo.EXPECT().AddUserSlugWithTtl(gomock.Any(), userID, slugID, name, ttl).Return(int64(0), nil)

		slugRepo := mock_modify_slug.NewMockslugRepo(ctrl)
		slugRepo.EXPECT().GetID(gomock.Any(), name).Return(slugID, nil)

		eventsRepo := mock_modify_slug.NewMockeventsRepo(ctrl)
		eventsRepo.EXPECT().AddEvent(gomock.Any(), userID, slugID, services.AddSlug)

		transactor := mock_modify_slug.NewMocktransactor(ctrl)
		transactor.EXPECT().RunInTx(gomock.Any(), gomock.Any()).DoAndReturn(
			func(ctx context.Context, f func(ctx context.Context) error) error {
				return f(ctx)
			})

		s := modifyslug.New(zap.L(), slugRepo, usersRepo, eventsRepo, transactor)

		// Action.
		err := s.ModifySlugs(ctx, userID, adds, nil)

		// Assert.
		assert.NoError(t, err)
	})

	t.Run("success add and delete slugs", func(t *testing.T) {
		// Arrange.
		ctx := context.Background()
		userID := uuid.New()
		name := "add"
		ttl := time.Now()
		adds := []modifyslug.Slug{
			{
				Name: name,
				Ttl:  ttl,
			},
		}

		nameDelete := "delete"

		deleteSlugs := []string{nameDelete}

		slugID := int64(1)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		usersRepo := mock_modify_slug.NewMockusersRepo(ctrl)
		usersRepo.EXPECT().CreateUserIfNotExist(gomock.Any(), userID).Return(nil)
		usersRepo.EXPECT().AddUserSlugWithTtl(gomock.Any(), userID, slugID, name, ttl).Return(int64(0), nil)
		usersRepo.EXPECT().DeleteUserSlug(gomock.Any(), userID, slugID).Return(nil)

		slugRepo := mock_modify_slug.NewMockslugRepo(ctrl)
		slugRepo.EXPECT().GetID(gomock.Any(), name).Return(slugID, nil)
		slugRepo.EXPECT().GetID(gomock.Any(), nameDelete).Return(slugID, nil)

		eventsRepo := mock_modify_slug.NewMockeventsRepo(ctrl)
		eventsRepo.EXPECT().AddEvent(gomock.Any(), userID, slugID, services.AddSlug)
		eventsRepo.EXPECT().AddEvent(gomock.Any(), userID, slugID, services.DeleteSlug)

		transactor := mock_modify_slug.NewMocktransactor(ctrl)
		transactor.EXPECT().RunInTx(gomock.Any(), gomock.Any()).DoAndReturn(
			func(ctx context.Context, f func(ctx context.Context) error) error {
				return f(ctx)
			})

		s := modifyslug.New(zap.L(), slugRepo, usersRepo, eventsRepo, transactor)

		// Action.
		err := s.ModifySlugs(ctx, userID, adds, deleteSlugs)

		// Assert.
		assert.NoError(t, err)
	})
}
