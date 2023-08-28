package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	v1 "github.com/frutonanny/slug-service/internal/generated/server/v1"
	"github.com/frutonanny/slug-service/pkg/errcodes"
)

func (h *Handlers) PostDeleteSlug(eCtx echo.Context) error {
	ctx := eCtx.Request().Context()

	var req v1.DeleteSlugRequest
	if err := eCtx.Bind(&req); err != nil {
		return eCtx.JSON(http.StatusBadRequest, v1.DeleteSlugResponse{
			Error: &v1.Error{
				Code:    errcodes.BadRequestError,
				Message: err.Error(),
			},
		})
	}

	if err := h.deleteSlugService.DeleteSlug(ctx, req.Name); err != nil {
		return eCtx.JSON(http.StatusInternalServerError, v1.DeleteSlugResponse{
			Error: &v1.Error{
				Code:    errcodes.InternalError,
				Message: "internal server error",
			},
		})
	}

	return eCtx.JSON(http.StatusOK, v1.DeleteSlugResponse{})
}
