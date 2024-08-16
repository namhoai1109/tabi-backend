package user

import (
	"context"
	"net/http"
	"tabi-booking/internal/model"

	"github.com/labstack/echo/v4"
	httpcore "github.com/namhoai1109/tabi/core/http"
)

type HTTP struct {
	service Service
}

type Service interface {
	GetSurvey(ctx context.Context, userID int) (*model.Survey, error)
}

func NewHTTP(service Service, echoGroup *echo.Group) {
	http := HTTP{service}

	echoGroup.GET("/:id/survey", http.list)

}

func (h *HTTP) list(c echo.Context) error {
	userID, err := httpcore.ReqID(c)
	if err != nil {
		return err
	}

	resp, err := h.service.GetSurvey(c.Request().Context(), userID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, resp)
}
