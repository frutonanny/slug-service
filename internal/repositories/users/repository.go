package users

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/frutonanny/slug-service/internal/database"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type Repository struct {
	db *database.DB
}

func New(db *database.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateUserIfNotExist(ctx context.Context, userID uuid.UUID) error {
	const query = `insert into "users" (user_id) values($1) on conflict (user_id) do nothing;`

	_, err := r.db.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("exec insert: %v", err)
	}

	return nil
}

func (r *Repository) AddUserSlug(ctx context.Context, userID uuid.UUID, slugID int64, name string) (int64, error) {
	const query = `insert into "users_slugs" (user_id, slug_id, slug_name) values ($1, $2, $3) returning id;`

	var id int64
	if err := r.db.QueryRow(ctx, query, userID, slugID, name).Scan(&id); err != nil {
		return 0, fmt.Errorf("query row: %v", err)
	}

	return id, nil
}

func (r *Repository) AddUserSlugWithTtl(
	ctx context.Context,
	userID uuid.UUID,
	slugID int64,
	name string,
	ttl time.Time,
) (int64, error) {
	const query = `insert into "users_slugs" (user_id, slug_id, slug_name, slug_ttl) values ($1, $2, $3, $4) returning id;`

	var id int64
	if err := r.db.QueryRow(ctx, query, userID, slugID, name, ttl).Scan(&id); err != nil {
		return 0, fmt.Errorf("query row: %v", err)
	}

	return id, nil
}

func (r *Repository) DeleteUserSlug(ctx context.Context, userID uuid.UUID, slugID int64) error {
	const query = `delete from "users_slugs" where "user_id" = $1 and "slug_id" = $2;`

	if _, err := r.db.Exec(ctx, query, userID, slugID); err != nil {
		return fmt.Errorf("exec insert: %v", err)
	}

	return nil
}

type Slug struct {
	ID        int64
	Name      string
	Ttl       time.Time
	DeletedAt time.Time
}

func (r *Repository) GetUserSlugs(ctx context.Context, userID uuid.UUID) ([]Slug, error) {
	const query = `select "slug_id", "slug_name", "slug_ttl", "deleted_at"
					from "users_slugs"
					left join "slugs" on "users_slugs"."slug_id" = "slugs"."id"
					where "users_slugs"."user_id" = $1;`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return []Slug{}, nil
		}

		return nil, fmt.Errorf("query slug name: %v", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	var result []Slug

	for rows.Next() {
		var (
			id        int64
			name      string
			ttl       sql.NullTime
			deletedAt sql.NullTime
		)

		if err := rows.Scan(&id, &name, &ttl, &deletedAt); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}

		slug := Slug{
			ID:   id,
			Name: name,
		}

		if ttl.Valid {
			slug.Ttl = ttl.Time
		}

		if deletedAt.Valid {
			slug.DeletedAt = deletedAt.Time
		}

		result = append(result, slug)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("rows: %w", err)
	}

	return result, nil
}

func (r *Repository) GetNextUser(ctx context.Context, previousID int64) (User, error) {
	const query = `select "id", "user_id" 
					from "users"
					where "id" > $1
  					order by id
  					limit 1;`

	var u User
	if err := r.db.QueryRow(ctx, query, previousID).Scan(&u.ID, &u.UserID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return u, ErrUserNotFound
		}

		return u, fmt.Errorf("scan user: %v", err)
	}

	return u, nil
}
