package me

import (
	"tabi-booking/internal/model"

	"github.com/labstack/echo/v4"
)

type HTTP struct {
	service Service
	autho   model.Autho
}

type Service interface {
	View(c echo.Context, authoUser *model.AuthoUser) (*model.UserResponse, error)
}

func NewHTTP(service Service, autho model.Autho, echoGroup *echo.Group) {
	http := HTTP{service, autho}

	// swagger:operation GET /user/me user-me UserMeView
	// ---
	// summary: View user information
	// responses:
	//   "200":
	//     description: user information
	//     schema:
	//       "$ref": "#/definitions/UserResponse"
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
	resp, err := h.service.View(c, h.autho.User(c))
	if err != nil {
		return err
	}

	return c.JSON(200, resp)
}
