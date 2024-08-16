package company

import (
	"context"
	"fmt"
	"tabi-booking/internal/model"

	"github.com/namhoai1109/tabi/core/logger"
	"github.com/namhoai1109/tabi/core/server"

	"gorm.io/gorm"
)

func (c *Company) checkExist(db *gorm.DB, ctx context.Context, id int) error {
	exist, err := c.companyDB.Exist(db, `representative_id = ?`, id)
	if err != nil {
		logger.LogError(ctx, fmt.Sprintf("Failed to check company exist: %v", err))
		return server.NewHTTPInternalError(fmt.Sprintf("Failed to check company exist: %v", err))
	}

	if !exist {
		logger.LogWarn(ctx, fmt.Sprintf("Company of representative_id %d does not exist", id))
		return server.NewHTTPValidationError(fmt.Sprintf("Company with representative_id %d does not exist", id))
	}

	return nil
}

func (c *Company) getBranchIDs(db *gorm.DB, ctx context.Context, rpID int) ([]int, error) {
	company := &model.Company{}
	if err := c.companyDB.View(db.Preload("Branches"), &company, `representative_id = ?`, rpID); err != nil {
		logger.LogError(ctx, fmt.Sprintf("Failed to view company: %v", err))
		return nil, server.NewHTTPInternalError(fmt.Sprintf("Failed to view company: %v", err))
	}

	branchIDs := []int{}
	for _, branch := range company.Branches {
		if branch.IsActive {
			branchIDs = append(branchIDs, branch.ID)
		}
	}

	return branchIDs, nil
}
