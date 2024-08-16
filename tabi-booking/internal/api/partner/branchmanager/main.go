package branchmanager

import (
	"context"
	"fmt"
	"tabi-booking/internal/model"
	"tabi-booking/internal/util/crypter"

	"github.com/namhoai1109/tabi/core/logger"
	"github.com/namhoai1109/tabi/core/server"
	"gorm.io/gorm"
)

func (s *BranchManagerService) Register(ctx context.Context, authoPartner *model.AuthoPartner, regData BranchManagerRegisterRequest) (*model.BranchManagerResponse, error) {
	if authoPartner.Role != model.AccountRoleRepresentative {
		return nil, server.NewHTTPAuthorizationError("Only representative can register branch manager")
	}

	if err := s.checkExistAccount(ctx, regData.Username, regData.Email, regData.Phone); err != nil {
		return nil, err
	}

	// Check if current branch already has branch manager
	if err := s.checkExistedBranchManager(ctx, regData.BranchID); err != nil {
		return nil, err
	}

	// Check if current branch belongs to current representative
	if err := s.checkExistBranch(ctx, authoPartner.ID, regData.BranchID); err != nil {
		return nil, err
	}

	branchManager := &model.BranchManager{}
	transErr := s.db.Transaction(func(tx *gorm.DB) error {
		// Check if branch manager already exist in branch

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
			Role:     model.AccountRoleBranchManager,
		}
		if err := s.accountDB.Create(tx, &account); err != nil {
			logger.LogError(ctx, fmt.Sprintf("Error when create account: %v", err))
			return server.NewHTTPInternalError("Error when create account")
		}

		branchManager = &model.BranchManager{
			Name:             regData.Name,
			AccountID:        account.ID,
			RepresentativeID: authoPartner.ID,
		}
		if err := s.branchManagerDB.Create(tx, &branchManager); err != nil {
			logger.LogError(ctx, fmt.Sprintf("Error when create branch manager: %v", err))
			return server.NewHTTPInternalError("Error when create branch manager")
		}

		updatedBranchData := &model.Branch{
			BranchManagerID: branchManager.ID,
		}
		if err := s.branchDB.Update(tx, updatedBranchData, `id = ?`, regData.BranchID); err != nil {
			logger.LogError(ctx, fmt.Sprintf("Error when update branch: %v", err))
			return server.NewHTTPInternalError("Error when update branch")
		}

		// TODO: The method below should be extracted to a new handler
		branchManager.Account = account
		return nil
	})

	if transErr != nil {
		return nil, transErr
	}

	return branchManager.ToResponse(), nil
}
