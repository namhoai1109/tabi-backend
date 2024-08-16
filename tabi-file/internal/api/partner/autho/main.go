package autho

import (
	"tabi-file/internal/model"

	"github.com/labstack/echo/v4"
)

// Admin returns admin data stored in jwt token
func (s *File) Partner(c echo.Context) *model.AuthoPartner {
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
