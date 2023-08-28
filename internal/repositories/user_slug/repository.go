package user_slug

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

func (r *Repository) AddUserSlug(ctx context.Context, user uuid.UUID, ids int64, name string) error {
	// todo
	return nil
}

func (r *Repository) DeleteUserSlug(ctx context.Context, userID uuid.UUID, ids int64, name string) error {
	// todo
	return nil
}

func (r *Repository) GetUserSlug(ctx context.Context, userID uuid.UUID) ([]string, error) {
	const query = `select "slug_name" from "users_slugs" where "user_id"= $1 and slug_ttl > now();`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("query: %v", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	var result []string

	for rows.Next() {
		var name string

		err = rows.Scan(&name)
		if err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		result = append(result, name)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("rows: %w", err)
	}

	return result, nil
}

func (r *Repository) DeleteUserSlugBySlugID(ctx context.Context, id int64) ([]uuid.UUID, error) {
	const query = `delete from "users_slugs" where "slug_id" = $1 returning "user_id";`

	rows, err := r.db.Query(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("query: %v", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	var result []uuid.UUID

	for rows.Next() {
		var name uuid.UUID

		err = rows.Scan(&name)
		if err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		result = append(result, name)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("rows: %w", err)
	}

	return result, nil
}

func (r *Repository) DeleteUserSlugByTtl(ctx context.Context, userID uuid.UUID) ([]int64, error) {
	const query = `select "slug_id" from "users_slugs" where "user_id"= $1 and slug_ttl > now();`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("query: %v", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	var result []int64

	for rows.Next() {
		var name int64

		err = rows.Scan(&name)
		if err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}
		result = append(result, name)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("rows: %w", err)
	}

	return result, nil
}
