package handlers

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	v1 "github.com/frutonanny/slug-service/internal/generated/server/v1"
	deleteslug "github.com/frutonanny/slug-service/internal/services/delete_slug"
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

		code := errcodes.InternalError
		msg := "internal server error"

		if errors.Is(err, deleteslug.ErrSlugNotFound) {
			code = errcodes.SlugNotFound
			msg = "internal server error"
		}

		return eCtx.JSON(http.StatusInternalServerError, v1.DeleteSlugResponse{
			Error: &v1.Error{
				Code:    code,
				Message: msg,
			},
		})
	}

	return eCtx.JSON(http.StatusOK, v1.DeleteSlugResponse{})
}
