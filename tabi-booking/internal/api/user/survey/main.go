package survey

import (
	"context"
	"tabi-booking/internal/model"

	"github.com/namhoai1109/tabi/core/server"
)

func (s *Survey) Create(ctx context.Context, authoUser *model.AuthoUser, req SurveyCreationRequest) (*model.Survey, error) {
	if authoUser.Role != model.AccountRoleClient {
		return nil, server.NewHTTPAuthorizationError("You are not allowed to access this resource")
	}

	survey := &model.Survey{
		UserID:     authoUser.ID,
		PlaceType:  req.PlaceType,
		Activities: req.Activities,
		Seasons:    req.Seasons,
	}

	if err := s.surveyDB.Create(s.db, survey); err != nil {
		return nil, server.NewHTTPInternalError("Error when creating survey")
	}

	return survey, nil
}
