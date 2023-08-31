//go:generate mockgen --source=service.go --destination=mock/service.go
package get_report

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"

	"github.com/frutonanny/wallet-service/internal/postgres"
	repoReport "github.com/frutonanny/wallet-service/internal/repositories/report"
	"github.com/frutonanny/wallet-service/pkg"
)

const (
	ReportsBucketName = "reports"
)

type logger interface {
	Info(msg string)
	Error(msg string)
}

type Repository interface {
	GetReport(ctx context.Context, period string) ([]repoReport.Service, error)
}

type MinioClient interface {
	PutObject(
		ctx context.Context,
		bucketName, objectName string,
		reader io.Reader,
		objectSize int64,
		opts minio.PutObjectOptions,
	) (info minio.UploadInfo, err error)
}

// dependencies умеет налету создавать репозиторий поверх *sql.DB, *sql.Tx.
// Нужен для написания юнит-тестов без подключения к базе.
type dependencies interface {
	NewRepository(db postgres.Database) Repository
}

type Service struct {
	logger         logger
	db             *sql.DB
	minioClient    MinioClient
	publicEndpoint string
	deps           dependencies
}

func New(logger logger, db *sql.DB, minioClient MinioClient, publicEndpoint string) *Service {
	return &Service{
		logger:         logger,
		db:             db,
		minioClient:    minioClient,
		publicEndpoint: publicEndpoint,

		deps: &dependenciesImpl{},
	}
}

func (s *Service) WithDependencies(deps dependencies) *Service {
	s.deps = deps
	return s
}

// GetReport отдает ссылку на CSV-файл, который лежит в хранилище minio. Файл содержит отчет за период period
// по всем услугам.
//
// - Получаем отчет из базы данных. В виде списка услуг за отчетный период period.
// - Преобразовываем полученный список в csv-файл в памяти.
// - Кладем преобразованный файл в minio-бакет отчетов.
// - Собираем ссылку на csv-файл.
func (s *Service) GetReport(ctx context.Context, period string) (string, error) {
	reportRepo := s.deps.NewRepository(s.db)

	// Получаем отчет из базы данных. В виде списка услуг за отчетный период period.
	report, err := reportRepo.GetReport(ctx, period)
	if err != nil {
		s.logger.Error(fmt.Sprintf("get report: %v", err))
		return "", fmt.Errorf("get report: %v", err)
	}

	// Преобразовываем полученный список в csv-файл в памяти.
	var b bytes.Buffer
	if err := writeToCsv(&b, report); err != nil {
		s.logger.Error(fmt.Sprintf("write to csv: %v", err))
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
		s.logger.Error(fmt.Sprintf("put object to minio: %v", err))
		return "", fmt.Errorf("put object to minio: %v", err)
	}

	return fmt.Sprintf("%s/%s/%s", s.publicEndpoint, ReportsBucketName, reportName), nil
}

// writeToCsv записывает полученный отчет в csv-файл.
func writeToCsv(wr io.Writer, report []repoReport.Service) error {
	csvWr := csv.NewWriter(wr)
	defer csvWr.Flush()

	for _, service := range report {
		record := []string{
			getServiceName(service.ServiceID),
			strconv.FormatInt(service.TotalRevenue, 10), // Общая выручка в копейках.
		}

		if err := csvWr.Write(record); err != nil {
			return fmt.Errorf("csv writer write: %v", err)
		}
	}

	return nil
}

func getServiceName(serviceID int64) string {
	if name, ok := pkg.Services[serviceID]; ok {
		return name
	}
	return fmt.Sprintf("Неизвестная услуга: %d", serviceID)
}
