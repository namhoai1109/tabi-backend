package branch

import "tabi-booking/internal/usecase/branch"

//swagger:model SaveBranchRequest
type SaveBranchRequest struct {
	Save *bool `json:"save" validate:"required"`
}

//swagger:model SaveBranchResponse
type SaveBranchResponse struct {
	Message string `json:"message"`
}

//swagger:model UserBranchListResponse
type BranchListResponse struct {
	Data  []*branch.PublicBranch `json:"data"`
	Total int                    `json:"total"`
}

type RatingBranchRequest struct {
	Rating    int    `json:"rating" validate:"required,min=1,max=5"`
	BookingID int    `json:"booking_id" validate:"required"`
	RoomID    int    `json:"room_id" validate:"required"`
	Comment   string `json:"comment"`
}
