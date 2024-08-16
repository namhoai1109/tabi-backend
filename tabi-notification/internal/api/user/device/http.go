package device

import (
	"context"
	"net/http"
	"tabi-notification/internal/model"

	"github.com/labstack/echo/v4"
)

type HTTP struct {
	service Service
	autho   model.Autho
}

type Service interface {
	Create(ctx context.Context, autho *model.AuthoUser, req DeviceCreationRequest) error
	Activate(ctx context.Context, autho *model.AuthoUser, req DeviceActivationRequest) error
	View(ctx context.Context, autho *model.AuthoUser, token string) (*model.Device, error)
}

func NewHTTP(service Service, autho model.Autho, echoGroup *echo.Group) {
	http := HTTP{service, autho}

	echoGroup.POST("", http.create)

	echoGroup.PATCH("/activation", http.activate)

	echoGroup.GET("/:pushToken", http.view)
}

func (h *HTTP) create(c echo.Context) error {
	req := DeviceCreationRequest{}
	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := h.service.Create(c.Request().Context(), h.autho.User(c), req); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *HTTP) view(c echo.Context) error {
	token := c.Param("pushToken")

	resp, err := h.service.View(c.Request().Context(), h.autho.User(c), token)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *HTTP) activate(c echo.Context) error {
	req := DeviceActivationRequest{}
	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := h.service.Activate(c.Request().Context(), h.autho.User(c), req); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}
