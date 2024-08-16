package user

import (
	"context"
	"fmt"
	"tabi-booking/internal/model"
	"tabi-booking/internal/util/crypter"
	structutil "tabi-booking/internal/util/struct"
	"time"

	"github.com/namhoai1109/tabi/core/logger"
	"github.com/namhoai1109/tabi/core/server"
)

func (s *AuthenUser) checkExistAccount(ctx context.Context, username, email, phone string) error {
	exist, err := s.accountDB.Exist(s.db, `(username = ? OR email = ? OR phone = ?)  AND role = ?`, username, email, phone, model.AccountRoleClient)
	if err != nil {
		logger.LogError(ctx, fmt.Sprintf("Error when check exist account: %v", err))
		return server.NewHTTPInternalError("Error when check exist account")
	}

	if exist {
		return server.NewHTTPValidationError("Account already exist")
	}

	return nil
}

func (s *AuthenUser) getToken(ctx context.Context, account *model.Account, userID int, remember bool) (*model.AuthToken, error) {
	userTokenClaims := &model.UserTokenClaims{
		ID:       userID,
		Username: account.Username,
		Email:    account.Email,
		Phone:    account.Phone,
		Role:     model.AccountRoleClient,
	}
	claims := structutil.ToMap(userTokenClaims)

	// default expire time is 24 hours
	timeExpire := time.Time.Add(time.Now(), 24*time.Hour)
	// if remember, expire time is 1 year
	if remember {
		timeExpire = time.Time.Add(time.Now(), 365*24*time.Hour)
	}
	token, expiresIn, err := s.jwt.GenerateToken(claims, &timeExpire)
	if err != nil {
		logger.LogError(ctx, fmt.Sprintf("Error when generate token: %v", err))
		return nil, server.NewHTTPInternalError("Error when generate token")
	}

	refreshToken := crypter.GenRefreshToken(s.cfg.JwtUserSecret)
	updates := map[string]interface{}{
		"refresh_token": refreshToken,
		"last_login":    time.Now(),
	}

	if err := s.accountDB.Update(s.db, updates, account.ID); err != nil {
		logger.LogError(ctx, fmt.Sprintf("Error when update account: %v", err))
		return nil, server.NewHTTPInternalError("Error when update account")
	}

	return &model.AuthToken{
		AccessToken:  token,
		TokenType:    "Bearer",
		ExpiresIn:    expiresIn,
		RefreshToken: refreshToken,
	}, nil
}
