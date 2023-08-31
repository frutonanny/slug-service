package outbox

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	outboxrepo "github.com/frutonanny/slug-service/internal/repositories/outbox"
)

type jobsRepository interface {
	CreateJob(ctx context.Context, name, data string) error
	FindJob(ctx context.Context) (outboxrepo.Job, error)
	ReserveJob(ctx context.Context, id int64, until time.Time) error
	DeleteJob(ctx context.Context, jobID int64) error
}

type transactor interface {
	RunInTx(ctx context.Context, f func(context.Context) error) error
}

const (
	workers    = 5
	idleTime   = 10 * time.Second
	reserveFor = 10 * time.Minute
)

type Service struct {
	jobs       map[string]Job
	outboxRepo jobsRepository
	transactor transactor
	log        *zap.Logger
}

func New(
	jobsRepo jobsRepository,
	transactor transactor,
	log *zap.Logger,
) *Service {
	return &Service{
		outboxRepo: jobsRepo,
		transactor: transactor,
		log:        log,
		jobs:       map[string]Job{},
	}
}

func (s *Service) RegisterJob(job Job) error {
	if _, ok := s.jobs[job.Name()]; ok {
		return fmt.Errorf("job %q already registered", job.Name())
	}

	s.jobs[job.Name()] = job
	return nil
}

func (s *Service) MustRegisterJob(job Job) {
	if err := s.RegisterJob(job); err != nil {
		panic(fmt.Errorf("register job: %v", err))
	}
}

func (s *Service) Run(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)

	for i := 0; i < workers; i++ {
		eg.Go(func() error {
			for {
				if err := s.processJobs(ctx); err != nil {
					if ctx.Err() != nil {
						return nil
					}
					return err
				}

				select {
				case <-ctx.Done():
					return nil
				case <-time.After(idleTime):
				}
			}
		})
	}

	return eg.Wait()
}

func (s *Service) processJobs(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		if err := s.findAndProcessJob(ctx); err != nil {
			if errors.Is(err, outboxrepo.ErrNoJobs) {
				return nil
			}
			return err
		}
	}
}

func (s *Service) findAndProcessJob(ctx context.Context) error {
	var job outboxrepo.Job
	var err error

	if err := s.transactor.RunInTx(ctx, func(ctx context.Context) error {
		job, err = s.outboxRepo.FindJob(ctx)
		if err != nil {
			return fmt.Errorf("find job: %w", err)
		}

		if err := s.outboxRepo.ReserveJob(ctx, job.ID, time.Now().Local().Add(reserveFor)); err != nil {
			return fmt.Errorf("reserve job: %v", err)
		}

		return nil
	}); err != nil {
		if errors.Is(err, outboxrepo.ErrNoJobs) {
			return nil
		}

		return fmt.Errorf("run in tx: %v", err)
	}

	j, ok := s.jobs[job.Name]
	if !ok {
		s.log.Error("job is not registered", zap.Error(err))
		return s.outboxRepo.DeleteJob(context.Background(), job.ID)
	}

	func() {
		ctx, cancel := context.WithTimeout(ctx, j.ExecutionTimeout())
		defer cancel()

		err = j.Handle(ctx, job.Data)
	}()

	if err != nil {
		s.log.Error("handle job error", zap.Error(err))
	}

	// Намеренно удаляем job с помощью context.Background(), чтобы избежать случая,
	// когда job обрабатывается, но ctx уже закрыт.
	if err := s.outboxRepo.DeleteJob(context.Background(), job.ID); err != nil {
		s.log.Error("delete job", zap.Error(err))
	}

	return nil
}

func (s *Service) Put(ctx context.Context, name, data string) error {
	return s.outboxRepo.CreateJob(ctx, name, data)
}
