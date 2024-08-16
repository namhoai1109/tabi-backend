package autho

import (
	"encoding/json"
	"tabi-booking/internal/model"
	structutil "tabi-booking/internal/util/struct"

	"github.com/labstack/echo/v4"
	"github.com/namhoai1109/tabi/core/logger"
)

// define to implement Autho interface, but not use
func (s *Autho) Partner(c echo.Context) *model.AuthoPartner {
	return &model.AuthoPartner{}
}

func (s *Autho) User(c echo.Context) *model.AuthoUser {
	// id, _ := c.Get("id").(int)
	// username, _ := c.Get("username").(string)
	// email, _ := c.Get("email").(string)
	// phone, _ := c.Get("phone").(string)
	// role, _ := c.Get("role").(string)

	// userTokenClaims := &model.UserTokenClaims{
	// 	ID:       id,
	// 	Username: username,
	// 	Email:    email,
	// 	Phone:    phone,
	// 	Role:     role,
	// }

	// return &model.AuthoUser{
	// 	ID:       userTokenClaims.ID,
	// 	Username: userTokenClaims.Username,
	// 	Email:    userTokenClaims.Email,
	// 	Phone:    userTokenClaims.Phone,
	// 	Role:     userTokenClaims.Role,
	// }

	// a bit lord :)
	ctx := c.Request().Context()
	tokenClaims := structutil.ToMap(&model.UserTokenClaims{})
	for k := range tokenClaims {
		tokenClaims[k] = c.Get(k)
	}

	tokenMarshal, err := json.Marshal(tokenClaims)
	if err != nil {
		logger.LogError(ctx, "Failed to marshal token claims")
		return nil
	}

	authoUser := &model.AuthoUser{}
	if err := json.Unmarshal(tokenMarshal, authoUser); err != nil {
		logger.LogError(ctx, "Failed to unmarshal token claims")
		return nil
	}

	return authoUser
}
