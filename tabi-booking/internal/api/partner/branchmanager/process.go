package branchmanager

import (
	"context"
	"fmt"
	"tabi-booking/internal/model"

	"github.com/namhoai1109/tabi/core/logger"
	"github.com/namhoai1109/tabi/core/server"
)

func (s *BranchManagerService) checkExistAccount(ctx context.Context, username, email, phone string) error {
	exist, err := s.accountDB.Exist(s.db, `(username = ? OR email = ? OR phone = ? ) AND role IN (?)`, username, email, phone, []string{
		model.AccountRoleBranchManager,
		model.AccountRoleRepresentative,
	})
	if err != nil {
		logger.LogError(ctx, fmt.Sprintf("Error when check exist account: %v", err))
		return server.NewHTTPInternalError("Error when check exist account")
	}

	if exist {
		return server.NewHTTPValidationError("Account already exist")
	}

	return nil
}

func (s *BranchManagerService) checkExistedBranchManager(ctx context.Context, branchID int) error {
	existed, err := s.branchDB.CheckExitedBranchManager(s.db, branchID)
	if err != nil {
		logger.LogError(ctx, fmt.Sprintf("Error when checking existed branch manager: %v", err))
		return server.NewHTTPInternalError(fmt.Sprintf("Error when checking existed branch manager, err: %v", err))
	}

	if existed {
		return server.NewHTTPValidationError("Branch manager already exist, please remove old branch manager before register new one")
	}

	return nil
}

func (s *BranchManagerService) checkExistBranch(ctx context.Context, rpID int, branchID int) error {
	existed, err := s.companyDB.CheckExistedBranchOfRepresentative(s.db, rpID, branchID)
	if err != nil {
		logger.LogError(ctx, fmt.Sprintf("Error when checking existed branch of representative: %v", err))
		return server.NewHTTPInternalError(fmt.Sprintf("Error when checking existed branch of representative, err: %v", err))
	}

	if !existed {
		return server.NewHTTPValidationError("Branch not found")
	}

	return nil
}
