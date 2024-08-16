package branch

import (
	"context"
	"fmt"
	"tabi-booking/internal/model"

	"github.com/namhoai1109/tabi/core/logger"
	"github.com/namhoai1109/tabi/core/rbac"
	"github.com/namhoai1109/tabi/core/server"

	"gorm.io/gorm"
)

func (s *Branch) toBranchListResponse(branches []*model.Branch, count int64) *BranchListResponse {
	list := []*CustomBranchResponse{}
	for _, branch := range branches {
		typeResp := &TypeResponse{}
		if branch.Type != nil {
			typeResp = &TypeResponse{
				ID:      branch.Type.ID,
				LabelVI: branch.Type.LabelVI,
				LabelEN: branch.Type.LabelEN,
			}
		}

		list = append(list, &CustomBranchResponse{
			ID:            branch.ID,
			BranchName:    branch.BranchName,
			Address:       branch.Address,
			ProvinceCity:  branch.ProvinceCity,
			District:      branch.District,
			Ward:          branch.Ward,
			FullAddress:   branch.FullAddress,
			ReceptionArea: branch.ReceptionArea,
			Description:   branch.Description,
			TypeResponse:  typeResp,
		})
	}

	return &BranchListResponse{
		Total: count,
		Data:  list,
	}
}

func (s *Branch) checkExistBranch(db *gorm.DB, ctx context.Context, authoPartner *model.AuthoPartner, branchID int) error {
	if authoPartner.Role == model.AccountRoleBranchManager || authoPartner.Role == model.AccountRoleHost {
		exist, err := s.branchDB.Exist(db, `id = ? AND branch_manager_id = ?`, branchID, authoPartner.ID)
		if err != nil {
			logger.LogError(ctx, fmt.Sprintf("Failed to check branch exist: %v", err))
			return server.NewHTTPInternalError(fmt.Sprintf("Failed to check branch exist: %v", err))
		}

		if !exist {
			logger.LogWarn(ctx, fmt.Sprintf("Branch with id %d does not exist", branchID))
			return server.NewHTTPValidationError(fmt.Sprintf("Branch with id %d does not exist", branchID))
		}
	} else { // authoPartner.Role == model.AccountRoleRepresentative
		company, err := s.companyDB.GetCompanyByRepID(db, authoPartner.ID)
		if err != nil {
			logger.LogError(ctx, fmt.Sprintf("Failed to get company: %v", err))
			return server.NewHTTPInternalError(fmt.Sprintf("Failed to get company: %v", err))
		}

		exist, err := s.branchDB.Exist(db, `id = ? AND company_id = ?`, branchID, company.ID)
		if err != nil {
			logger.LogError(ctx, fmt.Sprintf("Failed to check branch exist: %v", err))
			return server.NewHTTPInternalError(fmt.Sprintf("Failed to check branch exist: %v", err))
		}

		if !exist {
			logger.LogWarn(ctx, fmt.Sprintf("Branch with id %d does not exist", branchID))
			return server.NewHTTPValidationError(fmt.Sprintf("Branch with id %d does not exist", branchID))
		}
	}

	return nil
}

// enforce checks permission to perform the action
func (s *Branch) enforce(authPartner *model.AuthoPartner, action string) error {
	if !s.rbac.Enforce(authPartner.Role, model.ObjectBranch, action) {
		return rbac.ErrForbiddenAction
	}
	return nil
}
