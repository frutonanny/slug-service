package get_user_slug_test

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"

	usersrepo "github.com/frutonanny/slug-service/internal/repositories/users"
	getuserslug "github.com/frutonanny/slug-service/internal/services/get_user_slug"
	mock_get_user_slug "github.com/frutonanny/slug-service/internal/services/get_user_slug/mocks"
)

func TestService_GetUserSlug(t *testing.T) {
	t.Run("success get slugs", func(t *testing.T) {
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

		//assert.Eventually(t, func() bool {
		//	// Assert
		//	return usersRepo.DeleteUserSlugAssertExpectations(&testing.T{}) && dbMock.AssertExpectations(&testing.T{})
		//}, 100*time.Millisecond, 50*time.Millisecond)
	})
}
