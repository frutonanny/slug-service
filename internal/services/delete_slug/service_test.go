package delete_slug_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	slugrepo "github.com/frutonanny/slug-service/internal/repositories/slug"
	deleteslugservice "github.com/frutonanny/slug-service/internal/services/delete_slug"
	mock_delete_slug "github.com/frutonanny/slug-service/internal/services/delete_slug/mocks"
)

func TestService_DeleteSlug(t *testing.T) {
	t.Run("success delete", func(t *testing.T) {
		// Arrange.
		ctx := context.Background()
		name := "123"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		slugRepo := mock_delete_slug.NewMockslugRepo(ctrl)
		slugRepo.EXPECT().Delete(gomock.Any(), name).Return(nil)

		s := deleteslugservice.New(zap.L(), slugRepo)

		// Action.
		err := s.DeleteSlug(ctx, name)

		// Assert.
		assert.NoError(t, err)
	})

	t.Run("delete with error", func(t *testing.T) {
		// Arrange.
		ctx := context.Background()
		name := "123"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		slugRepo := mock_delete_slug.NewMockslugRepo(ctrl)
		slugRepo.EXPECT().Delete(gomock.Any(), name).Return(slugrepo.ErrRepoSlugNotFound)

		s := deleteslugservice.New(zap.L(), slugRepo)

		// Action.
		err := s.DeleteSlug(ctx, name)

		// Assert.
		assert.Error(t, err)
	})

}
