package branch

import (
	"context"
	"fmt"
	"tabi-booking/internal/model"

	"github.com/imdatngo/gowhere"
	dbcore "github.com/namhoai1109/tabi/core/db"
	"github.com/namhoai1109/tabi/core/logger"
	"github.com/namhoai1109/tabi/core/server"
)

func (s *Branch) checkExistBranch(ctx context.Context, branchID int) error {
	exist, err := s.branchDB.Exist(s.db, `id = ?`, branchID)
	if err != nil {
		logger.LogError(ctx, fmt.Sprintf("Error when check exist branch: %v", err))
		return server.NewHTTPInternalError("Error when check exist branch")
	}

	if !exist {
		return server.NewHTTPValidationError("Branch not found")
	}

	return nil
}

func (s *Branch) getRatings(ctx context.Context, branchID int) ([]*model.Rating, error) {
	lq := &dbcore.ListQueryCondition{
		Filter: gowhere.WithConfig(gowhere.Config{Strict: true}),
	}
	lq.Filter.And("branch_id = ?", branchID)
	ratings := []*model.Rating{}
	if err := s.ratingDB.List(s.db, &ratings, lq, nil); err != nil {
		logger.LogError(ctx, fmt.Sprintf("Error when get ratings: %v", err))
		return nil, server.NewHTTPInternalError("Error when get ratings")
	}

	return ratings, nil
}
