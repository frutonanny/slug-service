package get_report_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	eventsrepo "github.com/frutonanny/slug-service/internal/repositories/events"
	"github.com/frutonanny/slug-service/internal/services/get_report"
	mock_get_report "github.com/frutonanny/slug-service/internal/services/get_report/mock"
)

func TestService_GetReport(t *testing.T) {
	t.Run("get report success", func(t *testing.T) {
		ctx := context.Background()
		userID := uuid.New()
		publicEndpoint := "publicEndpoint"
		period := time.Now().Format("2006-01")

		from, to := parsePeriod(t, period)

		var events []eventsrepo.UserReportEvent

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// Arrange.
		eventsRepo := mock_get_report.NewMockeventsRepo(ctrl)
		eventsRepo.EXPECT().GetReport(gomock.Any(), userID, from, to).Return(events, nil)

		minioClient := mock_get_report.NewMockminioClient(ctrl)
		minioClient.
			EXPECT().
			PutObject(
				ctx,
				gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			).
			Return(minio.UploadInfo{}, nil)

		getUserSlugRepo := mock_get_report.NewMockgetUserSlugsService(ctrl)
		getUserSlugRepo.EXPECT().GetUserSlugs(gomock.Any(), userID, true).Return([]string{}, nil)

		service := get_report.New(zap.L(), getUserSlugRepo, eventsRepo, minioClient, publicEndpoint)

		// Action.
		_, err := service.GetReport(ctx, userID, period)

		// Assert.
		assert.NoError(t, err)
	})

	t.Run("failed: get slug error", func(t *testing.T) {
		ctx := context.Background()
		userID := uuid.New()
		publicEndpoint := "publicEndpoint"
		period := time.Now().Format("2006-01")

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		errExpected := errors.New("error")

		// Arrange.
		eventsRepo := mock_get_report.NewMockeventsRepo(ctrl)

		minioClient := mock_get_report.NewMockminioClient(ctrl)

		getUserSlugRepo := mock_get_report.NewMockgetUserSlugsService(ctrl)
		getUserSlugRepo.EXPECT().GetUserSlugs(gomock.Any(), userID, true).Return(nil, errExpected)

		service := get_report.New(zap.L(), getUserSlugRepo, eventsRepo, minioClient, publicEndpoint)

		// Action.
		_, err := service.GetReport(ctx, userID, period)

		// Assert.
		assert.Error(t, err)
	})

	t.Run("get report failed", func(t *testing.T) {
		ctx := context.Background()
		userID := uuid.New()
		publicEndpoint := "publicEndpoint"
		period := time.Now().Format("2006-01")
		from, to := parsePeriod(t, period)
		errExpected := errors.New("error")

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// Arrange.
		eventsRepo := mock_get_report.NewMockeventsRepo(ctrl)
		eventsRepo.EXPECT().GetReport(gomock.Any(), userID, from, to).Return(nil, errExpected)

		minioClient := mock_get_report.NewMockminioClient(ctrl)

		getUserSlugRepo := mock_get_report.NewMockgetUserSlugsService(ctrl)
		getUserSlugRepo.EXPECT().GetUserSlugs(gomock.Any(), userID, true).Return([]string{}, nil)

		service := get_report.New(zap.L(), getUserSlugRepo, eventsRepo, minioClient, publicEndpoint)

		// Action.
		_, err := service.GetReport(ctx, userID, period)

		// Assert.
		assert.Error(t, err)
	})
}

func parsePeriod(t *testing.T, period string) (time.Time, time.Time) {
	from, err := time.Parse("2006-01", period)
	require.NoError(t, err)

	to := from.AddDate(0, 1, 0)

	return from, to
}
