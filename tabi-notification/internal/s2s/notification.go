package s2s

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"tabi-notification/internal/model"

	"github.com/namhoai1109/tabi/core/logger"
	"github.com/namhoai1109/tabi/core/server"
)

func (s *S2S) GetDeviceList(ctx context.Context, req map[string]interface{}) ([]*model.Device, error) {
	params := url.Values{}
	if req != nil {
		f, err := json.Marshal(req)
		if err != nil {
			return nil, err
		}
		params.Add("f", string(f))
	}

	resp, err := s.s2s.Get(ctx, s.cfg.NotificationEndpoint, fmt.Sprintf(`/devices?%v`, params.Encode()))
	if err != nil {
		message := fmt.Sprintf("error getting device list: %v", err)
		logger.LogError(ctx, message)
		return nil, server.NewHTTPInternalError(message)
	}

	if err := s.s2s.BuildError(resp); err != nil {
		return nil, err
	}

	devices := []*model.Device{}
	if err := json.Unmarshal(resp.Body(), &devices); err != nil {
		message := fmt.Sprintf("error unmarshal device list: %v", err)
		logger.LogError(ctx, message)
		return nil, server.NewHTTPInternalError(message)
	}

	return devices, nil
}

func (s *S2S) GetScheduleList(ctx context.Context, req map[string]interface{}) ([]*model.Schedule, error) {
	params := url.Values{}
	if req != nil {
		f, err := json.Marshal(req)
		if err != nil {
			return nil, err
		}
		params.Add("f", string(f))
	}

	resp, err := s.s2s.Get(ctx, s.cfg.NotificationEndpoint, fmt.Sprintf(`/schedules?%v`, params.Encode()))
	if err != nil {
		message := fmt.Sprintf("error getting schedule list: %v", err)
		logger.LogError(ctx, message)
		return nil, server.NewHTTPInternalError(message)
	}

	if err := s.s2s.BuildError(resp); err != nil {
		return nil, err
	}

	schedules := []*model.Schedule{}
	if err := json.Unmarshal(resp.Body(), &schedules); err != nil {
		message := fmt.Sprintf("error unmarshal schedule list: %v", err)
		logger.LogError(ctx, message)
		return nil, server.NewHTTPInternalError(message)
	}

	return schedules, nil
}

func (s *S2S) MarkNotified(ctx context.Context, id int) error {
	resp, err := s.s2s.Patch(ctx, nil, s.cfg.NotificationEndpoint, fmt.Sprintf(`/schedules/%d/notified`, id))
	if err != nil {
		message := fmt.Sprintf("error marking schedule notified: %v", err)
		logger.LogError(ctx, message)
		return server.NewHTTPInternalError(message)
	}

	if err := s.s2s.BuildError(resp); err != nil {
		return err
	}

	return nil
}
