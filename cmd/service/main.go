package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	conf "github.com/frutonanny/slug-service/internal/config"
	"github.com/frutonanny/slug-service/internal/database"
	serverGen "github.com/frutonanny/slug-service/internal/generated/server/v1"
	logg "github.com/frutonanny/slug-service/internal/logger"
	"github.com/frutonanny/slug-service/internal/minio"
	"github.com/frutonanny/slug-service/internal/repositories/events"
	outboxrepo "github.com/frutonanny/slug-service/internal/repositories/outbox"
	slugrepo "github.com/frutonanny/slug-service/internal/repositories/slug"
	userslugrepo "github.com/frutonanny/slug-service/internal/repositories/users"
	createslugservice "github.com/frutonanny/slug-service/internal/services/create_slug"
	deleteslugservice "github.com/frutonanny/slug-service/internal/services/delete_slug"
	getreportservice "github.com/frutonanny/slug-service/internal/services/get_report"
	getuserslugservice "github.com/frutonanny/slug-service/internal/services/get_user_slug"
	modifyuserslugservice "github.com/frutonanny/slug-service/internal/services/modify_slug"
	outboxservice "github.com/frutonanny/slug-service/internal/services/outbox"
	percentslugjob "github.com/frutonanny/slug-service/internal/services/outbox/jobs/percent_slug"
	sortinghat "github.com/frutonanny/slug-service/internal/services/sorting_hat"
)

var configFile string

func init() {
	flag.StringVar(
		&configFile,
		"config",
		"config/config.local.json",
		"Path to configuration file",
	)
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("run: %v", err)
	}
}

func run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	flag.Parse()
	f := flag.Lookup(conf.Arg)
	if f == nil {
		return errors.New("config arg must be set")
	}

	config := conf.Must(f.Value.String())
	logger := logg.Must()

	// Postgres.
	db := database.Must(config.DB.DSN)
	defer func() {
		if err := db.Close(); err != nil {
			logger.Error("close db error: %s", zap.Error(err))
		}
	}()

	database.MustMigrate(db.DB)

	// Swagger.
	swagger, err := serverGen.GetSwagger()
	if err != nil {
		return fmt.Errorf("get swagger: %v", err)
	}

	// Minio.
	minioClient := minio.Must(
		config.Minio.Endpoint,
		config.Minio.AccessKeyID,
		config.Minio.SecretAccessKey,
	)

	// Repositories.
	slugRepository := slugrepo.New(db)
	outboxRepository := outboxrepo.New(db)
	usersRepository := userslugrepo.New(db)
	eventsRepository := events.New(db)

	// Services.
	sortingHatService := sortinghat.New()
	outboxService := outboxservice.New(outboxRepository, db, logger)
	outboxService.MustRegisterJob(percentslugjob.New(logger, sortingHatService, usersRepository, eventsRepository, db))

	createSlugService := createslugservice.New(logger, outboxService, slugRepository, db)
	deleteSlugService := deleteslugservice.New(logger, slugRepository)
	modifyUserSlugService := modifyuserslugservice.New(logger, slugRepository, usersRepository, eventsRepository, db)
	getUserSlugService := getuserslugservice.New(logger, usersRepository, eventsRepository, db)
	getReportService := getreportservice.New(
		logger,
		getUserSlugService,
		eventsRepository,
		minioClient,
		config.Minio.PublicEndpoint,
	)

	srv, err := initServer(
		net.JoinHostPort(config.Service.Host, config.Service.Port),
		swagger,
		createSlugService,
		deleteSlugService,
		modifyUserSlugService,
		getUserSlugService,
		getReportService,
	)
	if err != nil {
		return fmt.Errorf("init server: %v", err)
	}

	eg, ctx := errgroup.WithContext(ctx)

	// Run service.
	eg.Go(func() error {
		if err := srv.Run(ctx); err != nil {
			return fmt.Errorf("run server: %v", err)
		}

		return nil
	})

	// Run outbox.
	eg.Go(func() error {
		if err := outboxService.Run(ctx); err != nil {
			return fmt.Errorf("run outbox: %v", err)
		}

		return nil
	})

	if err := eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}

	return nil
}
