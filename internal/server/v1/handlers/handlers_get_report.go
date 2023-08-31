package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	v1 "github.com/frutonanny/slug-service/internal/generated/server/v1"
	"github.com/frutonanny/slug-service/pkg/errcodes"
)

func (h *Handlers) PostGetReport(eCtx echo.Context) error {
	ctx := eCtx.Request().Context()

	var req v1.GetReportRequest
	if err := eCtx.Bind(&req); err != nil {
		return eCtx.JSON(http.StatusBadRequest, v1.GetReportResponse{
			Error: &v1.Error{
				Code:    errcodes.BadRequestError,
				Message: err.Error(),
			},
		})
	}

	url, err := h.getReportService.GetReport(ctx, req.UserID, req.Period)
	if err != nil {
		return eCtx.JSON(http.StatusInternalServerError, v1.GetReportResponse{
			Error: &v1.Error{
				Code:    errcodes.InternalError,
				Message: "internal server error",
			},
		})
	}

	return eCtx.JSON(http.StatusOK, v1.GetReportResponse{
		Data: &v1.GetReportData{
			Url: url,
		},
	})
}
