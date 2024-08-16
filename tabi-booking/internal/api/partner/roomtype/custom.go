package roomtype

import (
	"tabi-booking/internal/model"
	"time"

	"github.com/lib/pq"
	httpcore "github.com/namhoai1109/tabi/core/http"
)

// swagger:model RoomTypeListResponse
type RoomTypeListResponse struct {
	Data  []*model.RoomTypeResponse `json:"data"`
	Total int64                     `json:"total"`
}

// swagger:model RoomTypeCreateRequest
type RoomTypeCreateRequest struct {
	TypeName         string        `json:"type_name" validate:"required"`
	CheckInTime      time.Time     `json:"check_in_time" validate:"required"`
	CheckOutTime     time.Time     `json:"check_out_time" validate:"required"`
	IncludeBreakfast bool          `json:"include_breakfast"`
	RoomFacilities   pq.Int64Array `json:"room_facilities" validate:"required"`
}

// swagger:model RoomTypeUpdateRequest
type RoomTypeUpdateRequest struct {
	TypeName         *string        `json:"type_name,omitempty"`
	CheckInTime      *time.Time     `json:"check_in_time,omitempty"`
	CheckOutTime     *time.Time     `json:"check_out_time,omitempty"`
	IncludeBreakfast *bool          `json:"include_breakfast,omitempty"`
	RoomFacilities   *pq.Int64Array `json:"room_facilities,omitempty"`
}

// swagger:model LinkRoomTypeRequest
type LinkRoomTypeRequest struct {
	RoomTypeID *int  `json:"room_type_id" validate:"required"`
	LinkStatus *bool `json:"link_status" validate:"required"`
}

// swagger:model LinkStatusResponse
type LinkStatusResponse struct {
	Status *string `json:"status"`
}

// swagger:model RoomTypeForBMResponse
type RoomTypeForBMResponse struct {
	Data []*model.RoomTypeResponse `json:"data"`
}

// swagger:parameters PartnerRoomTypeList
type ListRequestRoomTypeCustom struct {
	httpcore.ListRequest
	// id of branch
	// default: 1
	// required: true
	// in: path
	ID string `json:"id" query:"id"`
}
