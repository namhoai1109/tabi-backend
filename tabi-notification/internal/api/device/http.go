package device

import (
	"context"
	"net/http"
	"tabi-notification/internal/model"

	"github.com/labstack/echo/v4"
	dbcore "github.com/namhoai1109/tabi/core/db"
	httpcore "github.com/namhoai1109/tabi/core/http"
)

type HTTP struct {
	service Service
}

type Service interface {
	List(ctx context.Context, lq *dbcore.ListQueryCondition) ([]*model.Device, error)
}

func NewHTTP(service Service, echoGroup *echo.Group) {
	http := HTTP{service}

	echoGroup.GET("", http.list)
}

func (h *HTTP) list(c echo.Context) error {
	lq, err := httpcore.ReqListQuery(c)
	if err != nil {
		return err
	}

	resp, err := h.service.List(c.Request().Context(), lq)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}
