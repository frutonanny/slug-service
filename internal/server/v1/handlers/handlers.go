package handlers

import (
	"context"

	"github.com/google/uuid"

	createslugservice "github.com/frutonanny/slug-service/internal/services/create_slug"
	modifyslug "github.com/frutonanny/slug-service/internal/services/modify_slug"
)

type createSlugService interface {
	CreateSlug(ctx context.Context, name string, options createslugservice.Options) error
}

type deleteSlugService interface {
	DeleteSlug(ctx context.Context, name string) error
}

type modifyUserSlugService interface {
	ModifySlugs(ctx context.Context, userID uuid.UUID, addSlugs []modifyslug.Slug, delete []string) error
}

type getUserSlugService interface {
	GetUserSlugs(ctx context.Context, userID uuid.UUID, sync bool) ([]string, error)
}

type getReportService interface {
	GetReport(ctx context.Context, userID uuid.UUID, period string) (string, error)
}

type Handlers struct {
	createSlugService     createSlugService
	deleteSlugService     deleteSlugService
	modifyUserSlugService modifyUserSlugService
	getUserSlugService    getUserSlugService
	getReportService      getReportService
}

func NewHandlers(
	createSlugService createSlugService,
	deleteSlugService deleteSlugService,
	modifyUserSlugService modifyUserSlugService,
	getUserSlugService getUserSlugService,
	getReportService getReportService,
) *Handlers {
	return &Handlers{
		createSlugService:     createSlugService,
		deleteSlugService:     deleteSlugService,
		modifyUserSlugService: modifyUserSlugService,
		getUserSlugService:    getUserSlugService,
		getReportService:      getReportService,
	}
}
