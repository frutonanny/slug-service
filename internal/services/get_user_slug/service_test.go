package get_user_slug_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	usersrepo "github.com/frutonanny/slug-service/internal/repositories/users"
	getuserslug "github.com/frutonanny/slug-service/internal/services/get_user_slug"
	mock_get_user_slug "github.com/frutonanny/slug-service/internal/services/get_user_slug/mocks"
)

func TestService_GetUserSlug(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Arrange.
		ctx := context.Background()
		userID := uuid.New()
		id := int64(1)
		name := "123"

		slugs := []usersrepo.Slug{
			{
				ID:   id,
				Name: name,
			},
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		usersRepo := mock_get_user_slug.NewMockusersRepo(ctrl)
		usersRepo.EXPECT().GetUserSlugs(gomock.Any(), userID).Return(slugs, nil)

		eventsRepo := mock_get_user_slug.NewMockeventsRepo(ctrl)

		transactor := mock_get_user_slug.NewMocktransactor(ctrl)

		s := getuserslug.New(zap.L(), usersRepo, eventsRepo, transactor)

		// Action.
		result, err := s.GetUserSlug(ctx, userID)

		// Assert.
		require.NoError(t, err)
		assert.EqualValues(t, []string{name}, result)
	})

	t.Run("non-blocking delete slugs", func(t *testing.T) {
		// Arrange.
		ctx := context.Background()
		userID := uuid.New()

		slugs := []usersrepo.Slug{
			// Удаленный сегмент.
			{
				ID:        int64(1),
				Name:      "deleted_slug",
				DeletedAt: time.Now().Add(-1 * time.Minute),
			},
			// Истекший сегмент.
			{
				ID:   int64(2),
				Name: "outdated_slug",
				Ttl:  time.Now().Add(-1 * time.Minute),
			},
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		usersRepo := mock_get_user_slug.NewMockusersRepo(ctrl)
		usersRepo.EXPECT().GetUserSlugs(gomock.Any(), userID).Return(slugs, nil)
		first := false
		usersRepo.EXPECT().DeleteUserSlug(gomock.Any(), userID, slugs[0].ID).
			DoAndReturn(func(ctx context.Context, user uuid.UUID, slugID int64) error {
				first = true
				return nil
			})
		second := false
		usersRepo.EXPECT().DeleteUserSlug(gomock.Any(), userID, slugs[1].ID).
			DoAndReturn(func(ctx context.Context, user uuid.UUID, slugID int64) error {
				second = true
				return nil
			})

		eventsRepo := mock_get_user_slug.NewMockeventsRepo(ctrl)
		eventsRepo.EXPECT().AddEvent(gomock.Any(), userID, slugs[0].ID, gomock.Any()).Return(int64(0), nil)
		eventsRepo.EXPECT().AddEvent(gomock.Any(), userID, slugs[1].ID, gomock.Any()).Return(int64(0), nil)

		transactor := mock_get_user_slug.NewMocktransactor(ctrl)
		transactor.EXPECT().RunInTx(gomock.Any(), gomock.Any()).DoAndReturn(
			func(ctx context.Context, f func(ctx context.Context) error) error {
				return f(ctx)
			})

		s := getuserslug.New(zap.L(), usersRepo, eventsRepo, transactor)

		// Action.
		result, err := s.GetUserSlug(ctx, userID)

		// Assert.
		require.NoError(t, err)
		assert.EqualValues(t, []string{}, result)
		assert.Eventually(
			t,
			func() bool {
				return first && second
			},
			100*time.Millisecond,
			50*time.Millisecond,
		)
	})
}
