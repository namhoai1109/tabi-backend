package generaltype

import (
	"context"
	"fmt"
	"tabi-booking/internal/model"

	dbcore "github.com/namhoai1109/tabi/core/db"
	"github.com/namhoai1109/tabi/core/logger"
	"github.com/namhoai1109/tabi/core/server"

	"github.com/imdatngo/gowhere"
)

func (s *GeneralType) getTypeByOrder(ctx context.Context, typeList []string, lang string, order int) ([]*model.GeneralType, error) {
	lq := &dbcore.ListQueryCondition{
		Filter: gowhere.WithConfig(gowhere.Config{Strict: true}),
	}
	generalTypeList := []*model.GeneralType{}
	if lang == model.GeneralTypeLangVI {
		lq.Sort = []string{"label_vi ASC"}
	} else {
		lq.Sort = []string{"label_en ASC"}
	}
	lq.PerPage = 1000
	lq.Filter.And(`class IN (?) AND "order" = ?`, typeList, order)
	if err := s.generalTypeDB.List(s.db, &generalTypeList, lq, nil); err != nil {
		logger.LogError(ctx, fmt.Sprintf("Failed to get general type list: %s", err))
		return nil, server.NewHTTPInternalError(fmt.Sprintf("Failed to get general type list: %s", err))
	}

	return generalTypeList, nil
}
