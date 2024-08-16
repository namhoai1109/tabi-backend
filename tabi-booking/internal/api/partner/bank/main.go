package bank

import (
	"context"
	"fmt"
	"tabi-booking/internal/model"
	structutil "tabi-booking/internal/util/struct"

	dbcore "github.com/namhoai1109/tabi/core/db"
	"github.com/namhoai1109/tabi/core/logger"
	"github.com/namhoai1109/tabi/core/server"
	"gorm.io/gorm"
)

func (s *BankService) View(ctx context.Context, authoPartner *model.AuthoPartner, bankID int) (*CustomBankResponse, error) {
	if err := s.enforce(authoPartner, model.ActionView); err != nil {
		return nil, err
	}

	if err := s.checkBankAccountOwnership(ctx, authoPartner.ID, bankID); err != nil {
		return nil, err
	}

	bank := &model.Bank{}
	trxErr := s.db.Transaction(func(tx *gorm.DB) error {
		if err := s.bankDB.View(tx, &bank, `id = ?`, bankID); err != nil {
			logger.LogError(ctx, fmt.Sprintf("Failed to view bank account: %v", err))
			return server.NewHTTPInternalError(fmt.Sprintf("Failed to view bank account: %v", err))
		}

		return nil
	})

	if trxErr != nil {
		return nil, trxErr
	}

	return s.toResponse(bank), nil
}

func (s *BankService) Register(ctx context.Context, authoPartner *model.AuthoPartner, reqData *BankRegisterRequest, branchID int) (*CustomBankResponse, error) {
	if err := s.enforce(authoPartner, model.ActionCreate); err != nil {
		return nil, err
	}

	if err := s.checkValidAccount(ctx, reqData.AccountNumber); err != nil {
		return nil, err
	}
	if err := s.checkExistedAccount(ctx, reqData.AccountNumber); err != nil {
		return nil, err
	}
	if err := s.checkBranchOwnership(ctx, authoPartner.ID, branchID); err != nil {
		return nil, err
	}

	bank := &model.Bank{}
	trxErr := s.db.Transaction(func(tx *gorm.DB) error {
		bank = &model.Bank{
			BankID:        reqData.BankID,
			BankBranch:    reqData.BankBranch,
			AccountNumber: reqData.AccountNumber,
			AccountName:   reqData.AccountName,
			BranchID:      branchID,
		}

		if err := s.bankDB.Create(tx, &bank); err != nil {
			return err
		}

		return nil
	})

	if trxErr != nil {
		logger.LogError(ctx, fmt.Sprintf("Failed to register bank account: %v", trxErr))
		return nil, server.NewHTTPInternalError("Failed to register bank account")
	}

	return s.toResponse(bank), nil
}

func (s *BankService) List(ctx context.Context, authoPartner *model.AuthoPartner, branchID int, lq *dbcore.ListQueryCondition) (*BankListResponse, error) {
	if err := s.enforce(authoPartner, model.ActionViewAll); err != nil {
		return nil, err
	}
	if err := s.checkBranchOwnership(ctx, authoPartner.ID, branchID); err != nil {
		return nil, err
	}

	var count int64
	banks := []*model.Bank{}
	trxErr := s.db.Transaction(func(tx *gorm.DB) error {
		lq.Filter.And("branch_id = ?", branchID)
		if err := s.bankDB.List(tx, &banks, lq, &count); err != nil {
			logger.LogError(ctx, fmt.Sprintf("Failed to list bank account: %v", err))
			return server.NewHTTPInternalError(fmt.Sprintf("Failed to list bank account: %v", err))
		}

		return nil
	})

	if trxErr != nil {
		logger.LogError(ctx, fmt.Sprintf("Failed to list bank account: %v", trxErr))
		return nil, server.NewHTTPInternalError("Failed to list bank account")
	}

	return s.toListResponse(banks, count), nil
}

func (s *BankService) Update(ctx context.Context, authoPartner *model.AuthoPartner, bankID int, reqData *BankUpdateRequest) (*CustomBankResponse, error) {
	if err := s.enforce(authoPartner, model.ActionUpdate); err != nil {
		return nil, err
	}
	if err := s.checkBankAccountOwnership(ctx, authoPartner.ID, bankID); err != nil {
		return nil, err
	}

	updates := structutil.ToMap(*reqData)
	trxErr := s.db.Transaction(func(tx *gorm.DB) error {
		if err := s.bankDB.Update(tx, updates, bankID); err != nil {
			logger.LogError(ctx, fmt.Sprintf("Failed to update bank account: %v", err))
			return server.NewHTTPInternalError(fmt.Sprintf("Failed to update bank account: %v", err))
		}

		return nil
	})

	if trxErr != nil {
		return nil, trxErr
	}

	return s.View(ctx, authoPartner, bankID)
}

func (s *BankService) Delete(ctx context.Context, authoPartner *model.AuthoPartner, bankID int) error {
	if err := s.enforce(authoPartner, model.ActionDelete); err != nil {
		return err
	}
	if err := s.checkBankAccountOwnership(ctx, authoPartner.ID, bankID); err != nil {
		return err
	}

	trxErr := s.db.Transaction(func(tx *gorm.DB) error {
		if err := s.bankDB.Delete(tx, bankID); err != nil {
			return err
		}

		return nil
	})

	if trxErr != nil {
		logger.LogError(ctx, fmt.Sprintf("Failed to delete bank account: %v", trxErr))
		return server.NewHTTPInternalError("Failed to delete bank account")
	}

	return nil
}
