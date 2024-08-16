package me

import (
	"fmt"
	"tabi-booking/internal/model"

	"github.com/labstack/echo/v4"
	"github.com/namhoai1109/tabi/core/logger"
	"github.com/namhoai1109/tabi/core/server"
)

func (s *Me) View(c echo.Context, authoUser *model.AuthoUser) (*model.UserResponse, error) {
	if authoUser.Role != model.AccountRoleClient {
		logger.LogInfo(c.Request().Context(), "You don't have permission to view this resource")
		return nil, server.NewHTTPAuthorizationError("You don't have permission to view this resource")
	}

	user := &model.User{}
	if err := s.userDB.View(s.db.Preload("Account"), &user, `id = ?`, authoUser.ID); err != nil {
		logger.LogError(c.Request().Context(), fmt.Sprintf("Failed to view user: %v", err))
		return nil, server.NewHTTPInternalError(fmt.Sprintf("Failed to view user: %v", err))
	}

	return user.ToResponse(), nil
}
