package outbox

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/frutonanny/slug-service/internal/database"
)

var ErrNoJobs = errors.New("no jobs found")

type Job struct {
	ID   string
	Name string
	Data string
}

type Repository struct {
	db *database.DB
}

func New(db *database.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindJob(ctx context.Context) (Job, error) {
	const query = `select "id", "name", "data" 
					from "outbox" 
  					where "reserved_until" <= now() 
  					limit 1 for update skip locked;`

	var j Job
	if err := r.db.QueryRow(ctx, query).Scan(&j.ID, &j.Name, &j.Data); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Job{}, ErrNoJobs
		}
		return j, fmt.Errorf("scan job id: %v", err)
	}

	return j, nil
}

func (r *Repository) ReserveJob(ctx context.Context, id string, until time.Time) error {
	const query = `update "outbox" 
 					set "reserved_until" = $1 
					where "id" = $2;`

	_, err := r.db.Exec(ctx, query, until, id)
	if err != nil {
		return fmt.Errorf("exec update: %v", err)
	}

	return nil
}

func (r *Repository) CreateJob(ctx context.Context, name, data string) error {
	const query = `insert into "outbox" (name, data) values($1, $2);`

	_, err := r.db.Exec(ctx, query, name, data)
	if err != nil {
		return fmt.Errorf("exec insert: %v", err)
	}

	return nil
}

func (r *Repository) DeleteJob(ctx context.Context, jobID string) error {
	const query = `delete from "outbox" where id = $1;`

	_, err := r.db.Exec(ctx, query, jobID)
	if err != nil {
		return fmt.Errorf("exec delete: %v", err)
	}

	return nil
}
