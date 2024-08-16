package device

import (
	"context"
	"fmt"
	"tabi-notification/internal/model"

	"github.com/namhoai1109/tabi/core/logger"
	"github.com/namhoai1109/tabi/core/server"
)

func (s *Device) Create(ctx context.Context, autho *model.AuthoUser, req DeviceCreationRequest) error {
	exist, err := s.deviceDB.Exist(s.db, `push_token = ? AND user_id = ?`, req.PushToken, autho.ID)
	if err != nil {
		logger.LogError(ctx, fmt.Sprintf("error checking device existence: %v", err))
		return server.NewHTTPInternalError("error checking device existence")
	}

	if exist {
		return nil
	}

	device := &model.Device{
		PushToken: req.PushToken,
		UserID:    autho.ID,
		Brand:     req.Brand,
		Model:     req.Model,
		OS:        req.OS,
		OSVersion: req.OSVersion,
		IsActive:  true,
	}

	if err := s.deviceDB.Create(s.db, &device); err != nil {
		logger.LogError(ctx, fmt.Sprintf("error creating device: %v", err))
		return server.NewHTTPInternalError("error creating device")
	}

	return nil
}

func (s *Device) Activate(ctx context.Context, autho *model.AuthoUser, req DeviceActivationRequest) error {
	exist, err := s.deviceDB.Exist(s.db, `push_token = ? AND user_id = ?`, req.PushToken, autho.ID)
	if err != nil {
		logger.LogError(ctx, fmt.Sprintf("error checking device existence: %v", err))
		return server.NewHTTPInternalError("error checking device existence")
	}

	if !exist {
		return nil
	}

	if err := s.deviceDB.Update(s.db, map[string]interface{}{
		"is_active": *req.IsActive,
	}, `push_token = ? AND user_id = ?`, req.PushToken, autho.ID); err != nil {
		logger.LogError(ctx, fmt.Sprintf("error activating device: %v", err))
		return server.NewHTTPInternalError("error activating device")
	}

	return nil
}

func (s *Device) View(ctx context.Context, autho *model.AuthoUser, token string) (*model.Device, error) {
	if autho.Role != model.ClientRole {
		return nil, server.NewHTTPAuthorizationError("You are not authorized to access this resource")
	}

	device := model.Device{}
	if err := s.deviceDB.View(s.db, &device, `user_id = ? AND push_token = ?`, autho.ID, token); err != nil {
		logger.LogError(ctx, fmt.Sprintf("error getting device: %v", err))
		return nil, server.NewHTTPInternalError("error getting device")
	}

	return &device, nil
}
