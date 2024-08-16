package partner

import (
	"context"
	"fmt"
	"net/http"
	"tabi-booking/internal/model"
	"tabi-booking/internal/usecase/branch"
	"tabi-booking/internal/util/crypter"

	"github.com/namhoai1109/tabi/core/logger"
	"github.com/namhoai1109/tabi/core/server"

	"gorm.io/gorm"
)

var (
	ErrInvalidRefreshToken = server.NewHTTPError(http.StatusForbidden, "VALIDATION", "Refresh token is not valid.")
	ErrAccountNotExist     = server.NewHTTPValidationError("Account not exist")
)

func (s *AuthenPartner) RpRegister(ctx context.Context, regData RpRegistrationReq) (*model.AuthToken, error) {
	if err := s.checkExistAccount(ctx, regData.Username, regData.Email, regData.Phone); err != nil {
		return nil, err
	}

	transErr := s.db.Transaction(func(db *gorm.DB) error {
		// create account
		hashPassword, err := crypter.HashPassword(regData.Password)
		if err != nil {
			logger.LogError(ctx, fmt.Sprintf("Error when hash password: %v", err))
			return server.NewHTTPInternalError("Error when hash password")
		}

		account := &model.Account{
			Username: regData.Username,
			Password: hashPassword,
			Email:    regData.Email,
			Phone:    regData.Phone,
			Role:     model.AccountRoleRepresentative,
		}

		if err := s.accountDB.Create(db, &account); err != nil {
			logger.LogError(ctx, fmt.Sprintf("Error when create account: %v", err))
			return server.NewHTTPInternalError("Error when create account")
		}

		// create representative
		representative := &model.Representative{
			Name:      regData.FullName,
			AccountID: account.ID,
		}

		if err := s.representativeDB.Create(db, &representative); err != nil {
			logger.LogError(ctx, fmt.Sprintf("Error when create representative: %v", err))
			return server.NewHTTPInternalError("Error when create representative")
		}

		// create company
		company := &model.Company{
			CompanyName:      regData.CompanyName,
			ShortName:        regData.ShortName,
			Description:      regData.Description,
			WebsiteURL:       regData.WebsiteURL,
			TaxNumber:        regData.TaxNumber,
			RepresentativeID: representative.ID,
		}

		if err := s.companyDB.Create(db, &company); err != nil {
			logger.LogError(ctx, fmt.Sprintf("Error when create company: %v", err))
			return server.NewHTTPInternalError("Error when create company")
		}

		return nil
	})

	if transErr != nil {
		return nil, transErr
	}

	return s.Login(ctx, CredentialsPartnerReq{
		Username: regData.Username,
		Password: regData.Password,
	})
}

func (s *AuthenPartner) HstRegister(ctx context.Context, regData HstRegistrationReq) (*HSTRegisterResponse, error) {
	if err := s.checkExistAccount(ctx, regData.Username, regData.Email, regData.Phone); err != nil {
		return nil, err
	}

	branchID := 0
	transErr := s.db.Transaction(func(tx *gorm.DB) error {
		// create account
		hashPassword, err := crypter.HashPassword(regData.Password)
		if err != nil {
			logger.LogError(ctx, fmt.Sprintf("Error when hash password: %v", err))
			return server.NewHTTPInternalError("Error when hash password")
		}

		account := &model.Account{
			Username: regData.Username,
			Password: hashPassword,
			Email:    regData.Email,
			Phone:    regData.Phone,
			Role:     model.AccountRoleHost,
		}

		if err := s.accountDB.Create(tx, &account); err != nil {
			logger.LogError(ctx, fmt.Sprintf("Error when create account: %v", err))
			return server.NewHTTPInternalError("Error when create account")
		}

		// create host
		host := &model.BranchManager{
			Name:      regData.FullName,
			AccountID: account.ID,
		}

		if err := s.branchManagerDB.Create(tx, &host); err != nil {
			logger.LogError(ctx, fmt.Sprintf("Error when create host: %v", err))
			return server.NewHTTPInternalError("Error when create host")
		}

		branchCreation := branch.BranchCreationRequest{
			BranchName:            regData.BranchName,
			Address:               regData.Address,
			ProvinceCity:          regData.ProvinceCity,
			District:              regData.District,
			Ward:                  regData.Ward,
			Latitude:              regData.Latitude,
			Longitude:             regData.Longitude,
			Description:           regData.Description,
			ReceptionArea:         regData.ReceptionArea,
			MainFacilities:        regData.MainFacilities,
			TypeID:                regData.TypeID,
			CancellationTimeUnit:  regData.CancellationTimeUnit,
			CancellationTimeValue: regData.CancellationTimeValue,
			GeneralPolicy:         regData.GeneralPolicy,
			WebsiteURL:            regData.WebsiteURL,
			TaxNumber:             regData.TaxNumber,
			BranchManagerID:       &host.ID,
		}

		if _, err := s.branchUseCase.CreateBranch(tx, ctx, branchCreation, nil); err != nil {
			logger.LogError(ctx, fmt.Sprintf("Error when create branch: %v", err))
			return server.NewHTTPInternalError("Error when create branch")
		}

		branch := &model.Branch{}
		if err := s.branchDB.View(tx, branch, `branch_manager_id = ?`, host.ID); err != nil {
			logger.LogError(ctx, fmt.Sprintf("Error when get branch: %v", err))
			return server.NewHTTPInternalError("Error when get branch")
		}

		branchID = branch.ID
		return nil
	})

	if transErr != nil {
		return nil, transErr
	}

	authToken, err := s.Login(ctx, CredentialsPartnerReq{
		Username: regData.Username,
		Password: regData.Password,
	})
	if err != nil {
		return nil, err
	}

	return &HSTRegisterResponse{
		BranchID:     branchID,
		AccessToken:  authToken.AccessToken,
		TokenType:    authToken.TokenType,
		ExpiresIn:    authToken.ExpiresIn,
		RefreshToken: authToken.RefreshToken,
	}, nil
}

func (s *AuthenPartner) Login(ctx context.Context, credentials CredentialsPartnerReq) (*model.AuthToken, error) {
	// check exist account
	account := &model.Account{}
	if err := s.accountDB.View(s.db, &account, `username = ? AND role IN (?)`, credentials.Username, model.AccountRolesForPartner); err != nil {
		logger.LogError(ctx, fmt.Sprintf("Error when get account: %v", err))
		return nil, ErrAccountNotExist
	}

	// check password
	if !crypter.CompareHashAndPassword(account.Password, credentials.Password) {
		logger.LogWarn(ctx, "Wrong password")
		return nil, server.NewHTTPValidationError("Wrong password")
	}

	partnerID, err := s.getPartnerID(ctx, account)
	if err != nil {
		return nil, err
	}

	return s.getToken(ctx, account, partnerID)
}

func (s *AuthenPartner) RefreshToken(ctx context.Context, reqData RefreshTokenReq) (*model.AuthToken, error) {
	account := &model.Account{}
	if err := s.accountDB.View(s.db, &account, `refresh_token = ?`, reqData.RefreshToken); err != nil {
		logger.LogError(ctx, fmt.Sprintf("Error when get account: %v", err))
		return nil, ErrAccountNotExist
	}

	if account == nil {
		logger.LogWarn(ctx, "Refresh token not exist")
		return nil, ErrInvalidRefreshToken
	}

	if !crypter.ValidateRefreshToken(reqData.RefreshToken, s.cfg.JwtPartnerSecret) {
		logger.LogWarn(ctx, "Refresh token invalid")
		return nil, ErrInvalidRefreshToken
	}

	partnerID, err := s.getPartnerID(ctx, account)
	if err != nil {
		return nil, err
	}

	return s.getToken(ctx, account, partnerID)
}

func (s *AuthenPartner) Delete(ctx context.Context, rpID int) error {
	representative, err := s.checkExistRPByID(ctx, rpID)
	if err != nil {
		return err
	}

	trxErr := s.db.Transaction(func(db *gorm.DB) error {

		if err := s.accountDB.Delete(db, representative.Account.ID); err != nil {
			logger.LogError(ctx, fmt.Sprintf("Error when delete account: %v", err))
			return server.NewHTTPInternalError("Error when delete account")
		}

		if err := s.representativeDB.Delete(db, representative.ID); err != nil {
			logger.LogError(ctx, fmt.Sprintf("Error when delete representative: %v", err))
			return server.NewHTTPInternalError("Error when delete representative")
		}

		if err := s.companyDB.Delete(db, `representative_id = ?`, representative.ID); err != nil {
			logger.LogError(ctx, fmt.Sprintf("Error when delete company: %v", err))
			return server.NewHTTPInternalError("Error when delete company")
		}

		return nil
	})

	return trxErr
}
