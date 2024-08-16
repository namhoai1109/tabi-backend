package user

import (
	"context"
	"fmt"
	"tabi-booking/internal/model"

	"github.com/namhoai1109/tabi/core/server"
)

func (s *User) GetSurvey(ctx context.Context, userID int) (*model.Survey, error) {
	exist, err := s.userDB.Exist(s.db, userID)
	if err != nil {
		return nil, server.NewHTTPInternalError(fmt.Sprintf("Error when checking exist user: %v", err))
	}

	if !exist {
		return nil, server.NewHTTPValidationError("User not found")
	}

	existSurvey, err := s.surveyDB.Exist(s.db, `user_id = ?`, userID)
	if err != nil {
		return nil, server.NewHTTPInternalError(fmt.Sprintf("Error when checking exist survey: %v", err))
	}

	if !existSurvey {
		return nil, server.NewHTTPValidationError("Survey not found")
	}

	survey := &model.Survey{}
	if err := s.surveyDB.View(s.db, &survey, `user_id = ?`, userID); err != nil {
		return nil, server.NewHTTPInternalError(fmt.Sprintf("Error when getting survey: %v", err))
	}

	return survey, nil
}
