package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	v1 "github.com/frutonanny/slug-service/internal/generated/server/v1"
	createslugservice "github.com/frutonanny/slug-service/internal/services/create_slug"
	"github.com/frutonanny/slug-service/pkg/errcodes"
)

func (h *Handlers) PostCreateSlug(eCtx echo.Context) error {
	ctx := eCtx.Request().Context()

	var req v1.CreateSlugRequest
	if err := eCtx.Bind(&req); err != nil {
		return eCtx.JSON(http.StatusBadRequest, v1.CreateSlugResponse{
			Error: &v1.Error{
				Code:    errcodes.BadRequestError,
				Message: err.Error(),
			},
		})
	}

	options := createslugservice.Options{}
	if req.Options != nil {
		if req.Options.Percent != nil {
			options.Percent = req.Options.Percent
		}
	}

	if err := h.createSlugService.CreateSlug(ctx, req.Name, options); err != nil {
		return eCtx.JSON(http.StatusInternalServerError, v1.CreateSlugResponse{
			Error: &v1.Error{
				Code:    errcodes.InternalError,
				Message: "internal server error",
			},
		})
	}

	return eCtx.JSON(http.StatusOK, v1.CreateSlugResponse{})
}
