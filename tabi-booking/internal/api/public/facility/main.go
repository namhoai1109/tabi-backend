package facility

import (
	"context"
	"fmt"
	"tabi-booking/internal/model"

	dbcore "github.com/namhoai1109/tabi/core/db"
	"github.com/namhoai1109/tabi/core/logger"
	"github.com/namhoai1109/tabi/core/server"
)

func (s *Facility) List(ctx context.Context, lq *dbcore.ListQueryCondition, lang string) ([]*FacilityResponse, error) {

	facilities := []*model.Facility{}
	if lang == model.FacilityLangVI {
		lq.Sort = []string{"class_vi ASC", "name_vi ASC"}
	} else {
		lq.Sort = []string{"class_en ASC", "name_en ASC"}
	}
	lq.PerPage = 1000
	if err := s.facilityDB.List(s.db, &facilities, lq, nil); err != nil {
		logger.LogError(ctx, fmt.Sprintf("Failed to list facilities: %v", err))
		return nil, server.NewHTTPInternalError(fmt.Sprintf("Failed to list facilities: %v", err))
	}
	facilitiesResp := s.transformToFacilityResponse(facilities, lang)

	return facilitiesResp, nil
}
