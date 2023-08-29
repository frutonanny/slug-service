package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	v1 "github.com/frutonanny/slug-service/internal/generated/server/v1"
	modifyslug "github.com/frutonanny/slug-service/internal/services/modify_slug"
	"github.com/frutonanny/slug-service/pkg/errcodes"
)

func (h *Handlers) PostModifyUserSlugs(eCtx echo.Context) error {
	ctx := eCtx.Request().Context()

	var req v1.ModifyUserSlugsRequest
	if err := eCtx.Bind(&req); err != nil {
		return eCtx.JSON(http.StatusBadRequest, v1.ModifyUserSlugsResponse{
			Error: &v1.Error{
				Code:    errcodes.BadRequestError,
				Message: err.Error(),
			},
		})
	}

	var addSlugs []modifyslug.Slug
	if req.Add != nil {
		for _, slug := range *req.Add {
			addSlug := modifyslug.Slug{
				Name: slug.Name,
			}

			if slug.Ttl != nil && !slug.Ttl.IsZero() {
				addSlug.Ttl = *slug.Ttl
			}

			addSlugs = append(addSlugs, addSlug)
		}
	}

	var deleteSlugs []string
	if req.Delete != nil {
		for _, slug := range *req.Delete {
			deleteSlugs = append(deleteSlugs, slug)
		}
	}

	if err := h.modifyUserSlugService.ModifySlugs(ctx, req.UserID, addSlugs, deleteSlugs); err != nil {
		return eCtx.JSON(http.StatusInternalServerError, v1.ModifyUserSlugsResponse{
			Error: &v1.Error{
				Code:    errcodes.InternalError,
				Message: "internal server error",
			},
		})
	}

	return eCtx.JSON(http.StatusOK, v1.ModifyUserSlugsResponse{})
}
