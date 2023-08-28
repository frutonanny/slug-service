package handlers

import (
	"context"

	createslugservice "github.com/frutonanny/slug-service/internal/services/create_slug"
	modifyslug "github.com/frutonanny/slug-service/internal/services/modify_slug"
	"github.com/google/uuid"
)

type createSlugService interface {
	CreateSlug(ctx context.Context, name string, options createslugservice.Options) error
}

type deleteSlugService interface {
	DeleteSlug(ctx context.Context, name string) error
}

type modifyUserSlugService interface {
	ModifySlugs(ctx context.Context, userID uuid.UUID, addSlugs []modifyslug.AddSlug, delete []string) error
}

type getUserSlugService interface {
	GetUserSlug(ctx context.Context, userID uuid.UUID) ([]string, error)
}

type Handlers struct {
	createSlugService     createSlugService
	deleteSlugService     deleteSlugService
	modifyUserSlugService modifyUserSlugService
	getUserSlugService    getUserSlugService
}

func NewHandlers(
	createSlugService createSlugService,
	deleteSlugService deleteSlugService,
	modifyUserSlugService modifyUserSlugService,
	getUserSlugService getUserSlugService,
) *Handlers {
	return &Handlers{
		createSlugService:     createSlugService,
		deleteSlugService:     deleteSlugService,
		modifyUserSlugService: modifyUserSlugService,
		getUserSlugService:    getUserSlugService,
	}
}
