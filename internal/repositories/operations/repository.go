package operations

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

func (r *Repository) AddOperation(ctx context.Context, userID uuid.UUID, slugID int64, event string) error {
	const query = `insert into "users_slugs_history" (user_id, slug_id, event) values($1, $2, $3);`

	_, err := r.db.Exec(ctx, query, userID, slugID, event)
	if err != nil {
		return fmt.Errorf("exec insert: %v", err)
	}

	return nil
}
