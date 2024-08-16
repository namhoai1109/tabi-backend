package company

import (
	"context"
	"fmt"
	"reflect"
	"tabi-booking/internal/model"
	structutil "tabi-booking/internal/util/struct"

	"github.com/namhoai1109/tabi/core/logger"
	"github.com/namhoai1109/tabi/core/server"

	"gorm.io/gorm"
)

func (c *Company) ViewCompany(
	ctx context.Context,
	authoPartner *model.AuthoPartner) (*model.CompanyResponse, error) {
	if authoPartner.Role != model.AccountRoleRepresentative {
		return nil, server.NewHTTPAuthorizationError("You are not allowed to view company")
	}

	rpID := authoPartner.ID
	company := model.Company{}

	// transaction to view data
	trxErr := c.db.Transaction(func(db *gorm.DB) error {
		if err := c.checkExist(db, ctx, rpID); err != nil {
			return err
		}

		tx := db.Preload("Representative").Preload("Representative.Account")

		if err := c.companyDB.View(
			tx,
			&company,
			`representative_id = ?`, rpID); err != nil {
			logger.LogError(ctx, fmt.Sprintf("Failed to view company: %v", err))
			return server.NewHTTPInternalError(fmt.Sprintf("Failed to view company: %v", err))
		}

		return nil
	})

	if trxErr != nil {
		return nil, trxErr
	}

	return company.ToCompanyResponse(), nil
}

func (c *Company) UpdateCompany(
	ctx context.Context,
	authoPartner *model.AuthoPartner,
	data *CompanyUpdateRequest) (*model.CompanyResponse, error) {
	if authoPartner.Role != model.AccountRoleRepresentative {
		return nil, server.NewHTTPAuthorizationError("You are not allowed to update company")
	}

	// check if data is empty
	if reflect.ValueOf(*data).IsZero() {
		return nil, server.NewHTTPValidationError("Data is empty")
	}

	rpID := authoPartner.ID

	// transaction to update data
	trxErr := c.db.Transaction(func(db *gorm.DB) error {
		if err := c.checkExist(db, ctx, rpID); err != nil {
			return err
		}

		update := structutil.ToMap(*data)
		if update["email"] != nil {
			rp := &model.Representative{}
			if err := c.representativeDB.View(db.Preload("Account"), rp, `id = ?`, rpID); err != nil {
				logger.LogError(ctx, fmt.Sprintf("Failed to view representative: %v", err))
				return server.NewHTTPInternalError(fmt.Sprintf("Failed to view representative: %v", err))
			}

			if err := c.accountDB.Update(db, map[string]interface{}{
				"email": update["email"],
			}, `id = ?`, rp.Account.ID); err != nil {
				logger.LogError(ctx, fmt.Sprintf("Failed to update account: %v", err))
				return server.NewHTTPInternalError(fmt.Sprintf("Failed to update account: %v", err))
			}
			delete(update, "email")
		}

		if err := c.companyDB.Update(db, update, `representative_id = ?`, rpID); err != nil {
			logger.LogError(ctx, fmt.Sprintf("Failed to update company: %v", err))
			return server.NewHTTPInternalError(fmt.Sprintf("Failed to update company: %v", err))
		}

		return nil
	})

	if trxErr != nil {
		return nil, trxErr
	}

	return c.ViewCompany(ctx, authoPartner)
}

func (s *Company) AnalyzeRevenues(ctx context.Context, authoPartner *model.AuthoPartner, year int) ([]*model.RevenueAnalysisData, error) {
	if authoPartner.Role != model.AccountRoleRepresentative {
		return nil, server.NewHTTPAuthorizationError("You are not allowed to update company")
	}

	branchIDs, err := s.getBranchIDs(s.db, ctx, authoPartner.ID)
	if err != nil {
		return nil, err
	}

	resp, err := s.branchDB.AnalyzeBranchesRevenue(s.db, branchIDs, year)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to analyze branches revenue: %v", err)
		logger.LogError(ctx, errMsg)
		return nil, server.NewHTTPInternalError(errMsg)
	}

	return resp, nil
}

func (s *Company) AnalyzeBookingRequestQuantity(ctx context.Context, authoPartner *model.AuthoPartner, year int) ([]*model.BookingRequestQuantityAnalysisData, error) {
	if authoPartner.Role != model.AccountRoleRepresentative {
		return nil, server.NewHTTPAuthorizationError("You are not allowed to update company")
	}

	branchIDs, err := s.getBranchIDs(s.db, ctx, authoPartner.ID)
	if err != nil {
		return nil, err
	}

	resp, err := s.branchDB.AnalyzeBookingRequestQuantity(s.db, branchIDs, year)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to analyze branches revenue: %v", err)
		logger.LogError(ctx, errMsg)
		return nil, server.NewHTTPInternalError(errMsg)
	}

	return resp, nil
}
