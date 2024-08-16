package autho

import (
	"encoding/json"
	"tabi-payment/internal/model"
	structutil "tabi-payment/internal/util/struct"

	"github.com/labstack/echo/v4"
	"github.com/namhoai1109/tabi/core/logger"
)

// define to implement Autho interface, but not use
func (s *Autho) Partner(c echo.Context) *model.AuthoPartner {
	return &model.AuthoPartner{}
}

func (s *Autho) User(c echo.Context) *model.AuthoUser {
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
