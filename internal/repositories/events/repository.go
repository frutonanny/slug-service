package events

import (
	"context"
	"fmt"

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
