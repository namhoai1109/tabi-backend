package room

import (
	"tabi-booking/internal/model"

	httpcore "github.com/namhoai1109/tabi/core/http"
)

// swagger:model CreateRoomRequest
type CreateRoomRequest struct {
	RoomTypeID   int     `json:"room_type_id" validate:"required"`
	BedTypeID    int     `json:"bed_type_id" validate:"required"`
	RoomName     string  `json:"room_name" validate:"required"`
	MaxOccupancy int     `json:"max_occupancy" validate:"required"`
	Width        float32 `json:"width" validate:"required"`
	Length       float32 `json:"length" validate:"required"`
	MaxPrice     float64 `json:"max_price" validate:"required"`
	Quantity     int     `json:"quantity" validate:"required,min=1"`

	// facture reduction
	OnlineMethod *float64 `json:"online_method" validate:"required"`
	OnCashMethod *float64 `json:"on_cash_method" validate:"required"`
	NormalDay    *float64 `json:"normal_day" validate:"required"`
	Holiday      *float64 `json:"holiday" validate:"required"`
	Weekend      *float64 `json:"weekend" validate:"required"`

	// reservation reduction
	ReservationReduction []*ReservationReductionRequest `json:"reservation_reduction" validate:"required"`
}

// swagger:model ReservationReductionRequest
type ReservationReductionRequest struct {
	Quantity  *float64 `json:"quantity" validate:"required"`
	TimeUnit  string   `json:"time_unit" validate:"required"`
	Reduction *float64 `json:"reduction" validate:"required"`
}

// swagger:model RoomListItemResponse
type RoomListItemResponse struct {
	ID       int     `json:"id"`
	RoomType string  `json:"room_type"`
	RoomName string  `json:"room_name"`
	Status   string  `json:"status"`
	Quantity int     `json:"quantity"`
	MaxPrice float64 `json:"max_price"`

	BranchName        string `json:"branch_name"`
	BranchManagerName string `json:"branch_manager_name"`
}

// swagger:model RoomListResponse
type RoomListResponse struct {
	Total int64                  `json:"total"`
	Data  []RoomListItemResponse `json:"data"`
}

type CustomRoomType struct {
	model.RoomType
	Facilities []*model.Facility `json:"facilities"`
}

// swagger:model ViewRoomResponse
type ViewRoomResponse struct {
	ID           int     `json:"id"`
	RoomName     string  `json:"room_name"`
	MaxOccupancy int     `json:"max_occupancy"`
	Status       string  `json:"status"`
	Width        float32 `json:"width"`
	Length       float32 `json:"length"`
	MaxPrice     float64 `json:"max_price"`
	Quantity     int     `json:"quantity"`

	BedType               *model.GeneralType            `json:"bed_type"`
	FactureReduction      *model.FactureReduction       `json:"facture_reduction"`
	ReservationReductions []*model.ReservationReduction `json:"reservation_reductions"`
	RoomType              CustomRoomType                `json:"room_type"`
}

//swagger:model ListBookingsResponse
type ListBookingsResponse struct {
	Data  []*model.BookingResponse `json:"data"`
	Total int64                    `json:"total"`
}

// swagger:parameters PartnerRoomList
type ListRequest struct {
	httpcore.ListRequest
}

type FactureReductionUpdate struct {
	OnlineMethod float64 `json:"online_method,omitempty"`
	OnCashMethod float64 `json:"on_cash_method,omitempty"`
	NormalDay    float64 `json:"normal_day,omitempty"`
	Holiday      float64 `json:"holiday,omitempty"`
	Weekend      float64 `json:"weekend,omitempty"`
}

// swagger:model UpdateRoomRequest
type UpdateRoomRequest struct {
	RoomName     string  `json:"room_name,omitempty"`
	Status       string  `json:"status,omitempty"`
	MaxOccupancy int     `json:"max_occupancy,omitempty"`
	Width        float32 `json:"width,omitempty"`
	Length       float32 `json:"length,omitempty"`
	MaxPrice     float64 `json:"max_price,omitempty"`
	Quantity     int     `json:"quantity,omitempty"`
	BedTypeID    int     `json:"bed_type_id,omitempty"`
	RoomTypeID   int     `json:"room_type_id,omitempty"`

	// facture reduction
	OnlineMethod *float64 `json:"online_method,omitempty"`
	OnCashMethod *float64 `json:"on_cash_method,omitempty"`
	NormalDay    *float64 `json:"normal_day,omitempty"`
	Holiday      *float64 `json:"holiday,omitempty"`
	Weekend      *float64 `json:"weekend,omitempty"`

	ReservationReduction []*ReservationReductionRequest `json:"reservation_reduction,omitempty"`
}

type PaymentInfoResponse struct {
	RoomName    string `json:"room_name"`
	PaypalPayee string `json:"paypal_payee"`
}
