//go:generate mockgen --source=service.go --destination=mock/service.go
package get_report

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"go.uber.org/zap"

	eventsrepo "github.com/frutonanny/slug-service/internal/repositories/events"
)

const (
	ReportsBucketName = "reports"
)

type getUserSlugsService interface {
	GetUserSlugs(ctx context.Context, userID uuid.UUID, sync bool) ([]string, error)
}

type eventsRepo interface {
	GetReport(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]eventsrepo.UserReportEvent, error)
}

type minioClient interface {
	PutObject(
		ctx context.Context,
		bucketName, objectName string,
		reader io.Reader,
		objectSize int64,
		opts minio.PutObjectOptions,
	) (info minio.UploadInfo, err error)
}

type Service struct {
	log *zap.Logger

	getUserSlugsService getUserSlugsService
	eventsRepo          eventsRepo
	minioClient         minioClient
	publicEndpoint      string
}

func New(
	log *zap.Logger,
	getUserSlugsService getUserSlugsService,
	eventsRepo eventsRepo,
	minioClient minioClient,
	publicEndpoint string,
) *Service {
	return &Service{
		log:                 log,
		getUserSlugsService: getUserSlugsService,
		eventsRepo:          eventsRepo,
		minioClient:         minioClient,
		publicEndpoint:      publicEndpoint,
	}
}

func (s *Service) GetReport(ctx context.Context, userID uuid.UUID, period string) (string, error) {
	s.log.Info("start")
	// Анализирует и форматирует строку. Возвращает первый и последний дни месяца переданного года.
	from, to, err := parsePeriod(period)
	if err != nil {
		s.log.Error("parse period", zap.Error(err))
		return "", fmt.Errorf("parse period: %v", err)
	}

	// Подготавливает информацию о пользователе. Метод проходит по всем сегментам пользователя,
	// проверяет и удаляет сегменты, которые были удалены ранее или настало время автоматического удаления.
	if _, err := s.getUserSlugsService.GetUserSlugs(ctx, userID, true); err != nil {
		s.log.Error("prepare user for report: get user slugs", zap.Error(err))
		return "", fmt.Errorf("prepare user for report: get user slugs: %v", err)
	}

	// Получаем отчет из базы данных.
	report, err := s.eventsRepo.GetReport(ctx, userID, from, to)
	if err != nil {
		s.log.Error("get report", zap.Error(err))
		return "", fmt.Errorf("get report: %v", err)
	}

	// Преобразовываем полученный список в csv-файл в памяти.
	var b bytes.Buffer
	if err := writeToCsv(&b, report); err != nil {
		s.log.Error("write to csv", zap.Error(err))
		return "", fmt.Errorf("write to csv: %v", err)
	}

	reportName := fmt.Sprintf("report-%s-%s.csv", period, uuid.New())

	// Кладем преобразованный файл в minio-бакет отчетов.
	if _, err := s.minioClient.PutObject(
		ctx,
		ReportsBucketName,
		reportName,
		&b,
		int64(b.Len()),
		minio.PutObjectOptions{ContentType: "text/csv"},
	); err != nil {
		s.log.Error("put object to minio", zap.Error(err))
		return "", fmt.Errorf("put object to minio: %v", err)
	}

	s.log.Info(fmt.Sprintf("%s/%s/%s", s.publicEndpoint, ReportsBucketName, reportName))

	return fmt.Sprintf("%s/%s/%s", s.publicEndpoint, ReportsBucketName, reportName), nil
}

func parsePeriod(period string) (from, to time.Time, err error) {
	from, err = time.Parse("2006-01", period)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("parse: %v", err)
	}

	to = from.AddDate(0, 1, 0)

	return from, to, nil
}

// writeToCsv записывает полученный отчет в csv-файл.
func writeToCsv(wr io.Writer, report []eventsrepo.UserReportEvent) error {
	csvWr := csv.NewWriter(wr)
	defer csvWr.Flush()

	for _, event := range report {
		record := []string{
			event.UserID.String(),
			event.SlugName,
			event.EventName,
			event.CreatedAt.Format(time.RFC3339),
		}

		if err := csvWr.Write(record); err != nil {
			return fmt.Errorf("csv writer write: %v", err)
		}
	}

	return nil
}
