package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	v1 "github.com/frutonanny/slug-service/internal/generated/server/v1"
	"github.com/frutonanny/slug-service/pkg/errcodes"
)

func (h *Handlers) PostGetUserSlugs(eCtx echo.Context) error {
	ctx := eCtx.Request().Context()

	var req v1.GetUserSlugsRequest
	if err := eCtx.Bind(&req); err != nil {
		return eCtx.JSON(http.StatusBadRequest, v1.GetUserSlugsResponse{
			Error: &v1.Error{
				Code:    errcodes.BadRequestError,
				Message: err.Error(),
			},
		})
	}

	slugs, err := h.getUserSlugService.GetUserSlug(ctx, req.UserID)
	if err != nil {
		return eCtx.JSON(http.StatusInternalServerError, v1.GetUserSlugsResponse{
			Error: &v1.Error{
				Code:    errcodes.InternalError,
				Message: "internal server error",
			},
		})
	}

	return eCtx.JSON(http.StatusOK, v1.GetUserSlugsResponse{
		Data: &v1.GetUserSlugsData{Slugs: slugs},
	})
}
