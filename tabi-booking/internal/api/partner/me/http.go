package me

import (
	"context"
	"net/http"
	"tabi-booking/internal/model"

	"github.com/labstack/echo/v4"
)

type HTTP struct {
	service Service
	autho   model.Autho
}

type Service interface {
	View(ctx context.Context, authoPartner *model.AuthoPartner) (*PartnerInfoResponse, error)
}

func NewHTTP(service Service, autho model.Autho, echoGroup *echo.Group) {
	http := HTTP{service, autho}

	// swagger:operation GET /partner/me partner-me PartnerMeView
	// ---
	// summary: View partner information
	// responses:
	//   "200":
	//     description: partner information
	//     schema:
	//       "$ref": "#/definitions/PartnerInfoResponse"
	//   "400":
	//     "$ref": "#/responses/errDetails"
	//   "401":
	//     "$ref": "#/responses/errDetails"
	//   "403":
	//     "$ref": "#/responses/errDetails"
	//   "500":
	//     "$ref": "#/responses/errDetails"
	echoGroup.GET("", http.getMe)
}

func (h *HTTP) getMe(c echo.Context) error {
	resp, err := h.service.View(c.Request().Context(), h.autho.Partner(c))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}
