package slug

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/frutonanny/slug-service/internal/database"
)

var ErrRepoSlugNotFound = errors.New("slug not found")

type Options struct {
	Percent *int `json:"percent,omitempty"`
}

type Repository struct {
	db *database.DB
}

func New(db *database.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, name string, options Options) (int64, error) {
	const query = `insert into "slugs" (name, options) values ($1, $2) returning id;`

	b, err := json.Marshal(options)
	if err != nil {
		return 0, fmt.Errorf("marshal options: %v", err)
	}

	var id int64
	if err := r.db.QueryRow(ctx, query, name, string(b)).Scan(&id); err != nil {
		return 0, fmt.Errorf("query row: %v", err)
	}

	return id, nil
}

// Delete - метод удаления slug. Он полностью не удаляет из базы, а заносит данные в поле deleted_at.
func (r *Repository) Delete(ctx context.Context, name string) error {
	const query = `update "slugs" set "deleted_at" = now() where "name" = $1;`

	res, err := r.db.Exec(ctx, query, name)
	if err != nil {
		return fmt.Errorf("exec: %v", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affect: %w", err)
	}

	if rowsAffected == 0 {
		return ErrRepoSlugNotFound
	}

	return nil
}

func (r *Repository) GetID(ctx context.Context, name string) (int64, error) {
	const query = `select id from "slugs" where "name" = $1;`

	var id int64
	if err := r.db.QueryRow(ctx, query, name).Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrRepoSlugNotFound
		}

		return 0, fmt.Errorf("query row: %v", err)
	}

	return id, nil
}
