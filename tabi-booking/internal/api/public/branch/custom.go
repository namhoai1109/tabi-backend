package branch

import (
	"tabi-booking/internal/model"
	"tabi-booking/internal/usecase/branch"

	httpcore "github.com/namhoai1109/tabi/core/http"
)

//swagger:model PublicBranchListResponse
type BranchListResponse struct {
	Data  []*branch.PublicBranch `json:"data"`
	Total int                    `json:"total"`
}

// swagger:model PublicRoomListResponse
type RoomListResponse struct {
	Data  []*model.PublicRoom `json:"data"`
	Total int64               `json:"total"`
}

// swagger:parameters PublicBranchListRooms
type ListRequestPublicBranchCustom struct {
	httpcore.ListRequest
	// id of branch
	// default: 1
	// required: true
	// in: path
	ID string `json:"id" query:"id"`
}

//swagger:model FeaturedDestinationListResponse
type FeaturedDestinationListResponse struct {
	Data []string `json:"data"`
}

//swagger:model FeaturedBranchListResponse
type FeaturedBranchListResponse struct {
	Data []*branch.PublicBranch `json:"data"`
}

type RecommendedBranchListRequest struct {
	UserID int `json:"user_id"`
}
