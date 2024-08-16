package autho

import (
	"tabi-booking/internal/model"

	"github.com/labstack/echo/v4"
)

// Admin returns admin data stored in jwt token
func (s *Autho) Partner(c echo.Context) *model.AuthoPartner {
	id, _ := c.Get("id").(float64)
	username, _ := c.Get("username").(string)
	email, _ := c.Get("email").(string)
	role, _ := c.Get("role").(string)
	return &model.AuthoPartner{
		ID:       int(id),
		Username: username,
		Email:    email,
		Role:     role,
	}
}

// define to implement Autho interface, but not use
func (s *Autho) User(c echo.Context) *model.AuthoUser {
	return &model.AuthoUser{}
}
