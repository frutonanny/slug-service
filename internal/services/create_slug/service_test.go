package create_slug_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	slugrepo "github.com/frutonanny/slug-service/internal/repositories/slug"
	createslugservice "github.com/frutonanny/slug-service/internal/services/create_slug"
	mock_create_slug "github.com/frutonanny/slug-service/internal/services/create_slug/mocks"
)

func TestService_CreateSlug(t *testing.T) {
	t.Run("success with options", func(t *testing.T) {
		// Arrange.
		ctx := context.Background()
		id := int64(1)
		name := "123"
		percent := 50
		options := createslugservice.Options{
			Percent: &percent,
		}

		optionRepo := slugrepo.Options{
			Percent: &percent,
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		outboxService := mock_create_slug.NewMockoutboxService(ctrl)
		outboxService.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

		slugRepo := mock_create_slug.NewMockslugRepo(ctrl)
		slugRepo.EXPECT().Create(gomock.Any(), name, optionRepo).Return(id, nil)

		transactor := mock_create_slug.NewMocktransactor(ctrl)
		transactor.EXPECT().RunInTx(gomock.Any(), gomock.Any()).DoAndReturn(
			func(ctx context.Context, f func(ctx context.Context) error) error {
				return f(ctx)
			})

		s := createslugservice.New(zap.L(), outboxService, slugRepo, transactor)

		// Action.
		err := s.CreateSlug(ctx, name, options)

		// Assert.
		assert.NoError(t, err)
	})

	t.Run("outbox error with options", func(t *testing.T) {
		// Arrange.
		ctx := context.Background()
		name := "123"
		id := int64(1)
		percent := 50
		options := createslugservice.Options{
			Percent: &percent,
		}

		optionRepo := slugrepo.Options{
			Percent: &percent,
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		errExpected := errors.New("error")

		outboxService := mock_create_slug.NewMockoutboxService(ctrl)
		outboxService.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any()).Return(errExpected)

		slugRepo := mock_create_slug.NewMockslugRepo(ctrl)
		slugRepo.EXPECT().Create(gomock.Any(), name, optionRepo).Return(id, nil)

		transactor := mock_create_slug.NewMocktransactor(ctrl)
		transactor.EXPECT().RunInTx(gomock.Any(), gomock.Any()).DoAndReturn(
			func(ctx context.Context, f func(ctx context.Context) error) error {
				return f(ctx)
			})

		s := createslugservice.New(zap.L(), outboxService, slugRepo, transactor)

		// Action.
		err := s.CreateSlug(ctx, name, options)

		// Assert.
		assert.Error(t, err)
	})

	t.Run("success without options", func(t *testing.T) {
		// Arrange.
		ctx := context.Background()
		id := int64(1)
		name := "123"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		slugRepo := mock_create_slug.NewMockslugRepo(ctrl)
		slugRepo.EXPECT().Create(gomock.Any(), name, gomock.Any()).Return(id, nil)

		outboxService := mock_create_slug.NewMockoutboxService(ctrl)
		transactor := mock_create_slug.NewMocktransactor(ctrl)

		s := createslugservice.New(zap.L(), outboxService, slugRepo, transactor)

		// Action.
		err := s.CreateSlug(ctx, name, createslugservice.Options{})

		// Assert.
		assert.NoError(t, err)
	})
}
