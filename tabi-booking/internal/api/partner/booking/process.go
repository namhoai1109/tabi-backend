package booking

import (
	"context"
	"fmt"
	"tabi-booking/internal/model"

	"github.com/namhoai1109/tabi/core/logger"
	"github.com/namhoai1109/tabi/core/rbac"
	"github.com/namhoai1109/tabi/core/server"
)

func (s *Booking) getRoomIDs(ctx context.Context, bmID int) ([]int, error) {
	branch := &model.Branch{}
	if err := s.branchDB.View(s.db.Preload("Rooms"), &branch, `branch_manager_id = ?`, bmID); err != nil {
		logger.LogError(ctx, fmt.Sprintf("Error when getting branch: %v", err))
		return nil, server.NewHTTPInternalError(fmt.Sprintf("Error when getting branch, err: %v", err))
	}
	roomIDs := []int{}
	for _, room := range branch.Rooms {
		roomIDs = append(roomIDs, room.ID)
	}

	return roomIDs, nil
}

// enforce checks permission to perform the action
func (s *Booking) enforce(authPartner *model.AuthoPartner, action string) error {
	if !s.rbac.Enforce(authPartner.Role, model.ObjectBooking, action) {
		return rbac.ErrForbiddenAction
	}
	return nil
}
