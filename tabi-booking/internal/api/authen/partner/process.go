package partner

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

var (
	ErrRepresentativeNotExist = server.NewHTTPInternalError("Representative not exist")
)

func (s *AuthenPartner) checkExistAccount(ctx context.Context, username, email, phone string) error {
	exist, err := s.accountDB.Exist(s.db, `(username = ? OR email = ? OR phone = ? ) AND role IN (?)`, username, email, phone, model.AccountRolesForPartner)
	if err != nil {
		logger.LogError(ctx, fmt.Sprintf("Error when check exist account: %v", err))
		return server.NewHTTPInternalError("Error when check exist account")
	}

	if exist {
		return server.NewHTTPValidationError("Account already exist")
	}

	return nil
}

func (s *AuthenPartner) checkExistRPByID(ctx context.Context, rpID int) (*model.Representative, error) {
	representative := &model.Representative{}
	if err := s.representativeDB.View(s.db.Preload("Account"), &representative, rpID); err != nil {
		logger.LogError(ctx, fmt.Sprintf("Error when view representative: %v", err))
		return nil, ErrRepresentativeNotExist
	}

	if representative.Account == nil {
		logger.LogWarn(ctx, "Representative account is not exist")
		return nil, ErrRepresentativeNotExist
	}

	return representative, nil
}

func (s *AuthenPartner) getToken(ctx context.Context, account *model.Account, partnerID int) (*model.AuthToken, error) {
	// generate token
	partnerTokenClaims := &model.PartnerTokenClaims{
		ID:       partnerID,
		Username: account.Username,
		Role:     account.Role,
		Email:    account.Email,
	}
	claims := structutil.ToMap(partnerTokenClaims)

	oneYear := time.Time.Add(time.Now(), time.Hour*24*365)
	// oneYear := time.Time.Add(time.Now(), 365*24*time.Hour)
	token, expiresIn, err := s.jwt.GenerateToken(claims, &oneYear)
	if err != nil {
		logger.LogError(ctx, fmt.Sprintf("Error when generate token: %v", err))
		return nil, server.NewHTTPInternalError("Error when generate token")
	}

	// generate refresh token and update account
	refreshToken := crypter.GenRefreshToken(s.cfg.JwtPartnerSecret)
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

func (s *AuthenPartner) getPartnerID(ctx context.Context, account *model.Account) (int, error) {
	partnerID := 0
	if account.Role == model.AccountRoleRepresentative {
		rp := &model.Representative{}
		if err := s.representativeDB.View(s.db, &rp, `account_id = ?`, account.ID); err != nil {
			logger.LogError(ctx, fmt.Sprintf("Error when get representative: %v", err))
			return 0, server.NewHTTPInternalError("Error when get representative")
		}
		partnerID = rp.ID
	} else if account.Role == model.AccountRoleBranchManager || account.Role == model.AccountRoleHost {
		bm := &model.BranchManager{}
		if err := s.branchManagerDB.View(s.db, &bm, `account_id = ?`, account.ID); err != nil {
			logger.LogError(ctx, fmt.Sprintf("Error when get branch manager: %v", err))
			return 0, server.NewHTTPInternalError("Error when get branch manager")
		}
		partnerID = bm.ID
	}

	return partnerID, nil
}
