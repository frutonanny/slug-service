package events

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/frutonanny/slug-service/internal/database"
)

type Repository struct {
	db *database.DB
}

func New(db *database.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) AddEvent(ctx context.Context, userID uuid.UUID, slugID int64, event string) (int64, error) {
	const query = `insert into "events" (user_id, slug_id, event) values ($1, $2, $3) returning id;`

	var id int64
	if err := r.db.QueryRow(ctx, query, userID, slugID, event).Scan(&id); err != nil {
		return 0, fmt.Errorf("query row: %v", err)
	}

	return id, nil
}

func (r *Repository) AddEventWithCreatedAt(
	ctx context.Context,
	userID uuid.UUID,
	slugID int64,
	event string,
	createdAt time.Time,
) (int64, error) {
	const query = `insert into "events" (user_id, slug_id, event, created_at) values ($1, $2, $3, $4) returning id;`

	var id int64
	if err := r.db.QueryRow(ctx, query, userID, slugID, event, createdAt).Scan(&id); err != nil {
		return 0, fmt.Errorf("query row: %v", err)
	}

	return id, nil
}

func (r *Repository) GetReport(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]UserReportEvent, error) {
	const query = `select "slugs"."name", 
						"events"."event", 
						"events"."created_at" 
					from "events"
					left join "slugs" on "events"."slug_id" = "slugs"."id"
					where "events"."user_id" = $1 and 
					      "events"."created_at" >= $2 and "events"."created_at" <= $3 
  					order by "events"."id";`

	rows, err := r.db.Query(ctx, query, userID, from, to)
	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return []UserReportEvent{}, nil
		}

		return nil, fmt.Errorf("query report: %v", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	var result []UserReportEvent

	for rows.Next() {
		event := UserReportEvent{
			UserID: userID,
		}

		if err := rows.Scan(&event.SlugName, &event.EventName, &event.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}

		result = append(result, event)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("rows: %w", err)
	}

	return result, nil
}
