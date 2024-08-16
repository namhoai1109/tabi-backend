package bank

import (
	"context"
	"fmt"
	"regexp"
	"tabi-booking/internal/model"

	"github.com/namhoai1109/tabi/core/logger"
	"github.com/namhoai1109/tabi/core/rbac"
	"github.com/namhoai1109/tabi/core/server"
)

// enforce checks permission to perform the action
func (s *BankService) enforce(authPartner *model.AuthoPartner, action string) error {
	if !s.rbac.Enforce(authPartner.Role, model.ObjectBank, action) {
		return rbac.ErrForbiddenAction
	}
	return nil
}

func (s *BankService) checkValidAccount(ctx context.Context, bankAccount string) error {
	regex, err := regexp.Compile(`^[1-9][0-9]{15}$`)
	if err != nil {
		logger.LogError(ctx, fmt.Sprintf("Failed to compile regex when validating bank account: %v", err))
		return server.NewHTTPInternalError("Failed to validate bank account")
	}

	if !regex.MatchString(bankAccount) {
		return server.NewHTTPValidationError("Invalid bank account")
	}

	return nil
}

func (s *BankService) checkExistedAccount(ctx context.Context, bankAccount string) error {
	existed, err := s.bankDB.Exist(s.db, `account_number = ?`, bankAccount)
	if err != nil {
		logger.LogError(ctx, fmt.Sprintf("Failed to check existed account: %v", err))
		return server.NewHTTPInternalError("Failed to check existed account")
	}

	if existed {
		return server.NewHTTPValidationError("Bank account already exist")
	}

	return nil
}

func (s *BankService) checkBranchOwnership(ctx context.Context, rpID int, branchID int) error {
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

// checkOwnerBankAccount checks if the given bank account belongs to the representative
func (s *BankService) checkBankAccountOwnership(ctx context.Context, rpID int, bankID int) error {
	existed, err := s.companyDB.CheckBankAccountOwnership(s.db, rpID, bankID)
	if err != nil {
		logger.LogError(ctx, fmt.Sprintf("Error when checking bank account ownership: %v", err))
		return server.NewHTTPInternalError(fmt.Sprintf("Error when checking bank account ownership, err: %v", err))
	}

	if !existed {
		return server.NewHTTPValidationError("Bank account ownership is invalid")
	}

	return nil
}

func (s *BankService) toResponse(bank *model.Bank) *CustomBankResponse {
	return &CustomBankResponse{
		ID:            bank.ID,
		BankID:        bank.BankID,
		BankBranch:    bank.BankBranch,
		AccountNumber: bank.AccountNumber,
		AccountName:   bank.AccountName,
		BranchID:      bank.BranchID,
	}
}

func (s *BankService) toListResponse(banks []*model.Bank, count int64) *BankListResponse {
	list := []*CustomBankResponse{}
	for _, bank := range banks {
		list = append(list, s.toResponse(bank))
	}
	return &BankListResponse{
		Total: count,
		Data:  list,
	}
}
