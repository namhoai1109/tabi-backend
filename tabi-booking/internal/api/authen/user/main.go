package user

import (
	"context"
	"fmt"
	"tabi-booking/internal/model"
	"tabi-booking/internal/util/crypter"

	"github.com/namhoai1109/tabi/core/logger"
	"github.com/namhoai1109/tabi/core/server"
	"gorm.io/gorm"
)

func (s *AuthenUser) Register(ctx context.Context, regData RegistrationUserReq) (*model.AuthToken, error) {
	if err := s.checkExistAccount(ctx, regData.Username, regData.Email, regData.Phone); err != nil {
		return nil, err
	}

	trxErr := s.db.Transaction(func(tx *gorm.DB) error {
		hashPassword, err := crypter.HashPassword(regData.Password)
		if err != nil {
			logger.LogError(ctx, fmt.Sprintf("Error when hash password: %v", err))
			return server.NewHTTPInternalError("Error when hash password")
		}
		account := &model.Account{
			Username: regData.Username,
			Password: hashPassword,
			Phone:    regData.Phone,
			Email:    regData.Email,
			Role:     model.AccountRoleClient,
		}

		if err := s.accountDB.Create(tx, &account); err != nil {
			logger.LogError(ctx, fmt.Sprintf("Error when create account: %v", err))
			return server.NewHTTPInternalError("Error when create account")
		}

		user := &model.User{
			FirstName:   regData.FirstName,
			LastName:    regData.LastName,
			DateOfBirth: &regData.DoB,
			AccountID:   account.ID,
		}

		if err := s.userDB.Create(tx, &user); err != nil {
			logger.LogError(ctx, fmt.Sprintf("Error when create user: %v", err))
			return server.NewHTTPInternalError("Error when create user")
		}

		return nil
	})

	if trxErr != nil {
		return nil, trxErr
	}

	return s.Login(ctx, CredentialsUserReq{
		Identity: regData.Phone,
		Password: regData.Password,
	})
}

func (s *AuthenUser) Login(ctx context.Context, credentials CredentialsUserReq) (*model.AuthToken, error) {
	account := &model.Account{}
	identity := credentials.Identity
	if err := s.accountDB.View(s.db, &account, `(email = ? OR phone = ?) AND role = ?`, identity, identity, model.AccountRoleClient); err != nil {
		logger.LogError(ctx, fmt.Sprintf("Error when get account: %v", err))
		return nil, server.NewHTTPValidationError("Incorrect username or password")
	}

	if !crypter.CompareHashAndPassword(account.Password, credentials.Password) {
		logger.LogWarn(ctx, fmt.Sprintf("Incorrect password for account: %v", account.ID))
		return nil, server.NewHTTPValidationError("Incorrect username or password")
	}

	user := &model.User{}
	if err := s.userDB.View(s.db, &user, `account_id = ?`, account.ID); err != nil {
		logger.LogError(ctx, fmt.Sprintf("Error when get user: %v", err))
		return nil, server.NewHTTPInternalError("Error when get user")
	}

	return s.getToken(ctx, account, user.ID, credentials.Remember)
}
