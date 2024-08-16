package me

import (
	"context"
	"fmt"
	"tabi-booking/internal/model"
	structutil "tabi-booking/internal/util/struct"

	"github.com/namhoai1109/tabi/core/logger"
	"github.com/namhoai1109/tabi/core/server"
)

func (s *Me) View(ctx context.Context, authoPartner *model.AuthoPartner) (*PartnerInfoResponse, error) {
	resp := map[string]interface{}{}
	placeID := 0
	if authoPartner.Role == model.AccountRoleRepresentative {
		representative := &model.Representative{}
		if err := s.representativeDB.View(s.db.Preload("Account"), &representative, `id = ?`, authoPartner.ID); err != nil {
			logger.LogError(ctx, fmt.Sprintf("Failed to view representative: %v", err))
			return nil, server.NewHTTPInternalError(fmt.Sprintf("Failed to view representative: %v", err))
		}
		resp = structutil.ToMap(representative.ToRepresentativeResponse())

		company := &model.Company{}
		if err := s.companyDB.View(s.db, &company, `representative_id = ?`, authoPartner.ID); err != nil {
			logger.LogError(ctx, fmt.Sprintf("Failed to view company: %v", err))
			return nil, server.NewHTTPInternalError(fmt.Sprintf("Failed to view company: %v", err))
		}
		placeID = company.ID
	} else if authoPartner.Role == model.AccountRoleBranchManager || authoPartner.Role == model.AccountRoleHost {
		branchManager := &model.BranchManager{}
		if err := s.branchManagerDB.View(s.db.Preload("Account"), &branchManager, `id = ?`, authoPartner.ID); err != nil {
			logger.LogError(ctx, fmt.Sprintf("Failed to view branch manager: %v", err))
			return nil, server.NewHTTPInternalError(fmt.Sprintf("Failed to view branch manager: %v", err))
		}
		resp = structutil.ToMap(branchManager.ToResponse())

		branch := &model.Branch{}
		if err := s.branchDB.View(s.db, &branch, `branch_manager_id = ?`, branchManager.ID); err != nil {
			logger.LogError(ctx, fmt.Sprintf("Failed to view branch: %v", err))
			return nil, server.NewHTTPInternalError(fmt.Sprintf("Failed to view branch: %v", err))
		}
		placeID = branch.ID
	}

	return &PartnerInfoResponse{
		ID:       resp["id"].(int),
		Name:     resp["name"].(string),
		Username: resp["username"].(string),
		Email:    resp["email"].(string),
		Phone:    resp["phone"].(string),
		Role:     resp["role"].(string),
		PlaceID:  placeID,
	}, nil
}
