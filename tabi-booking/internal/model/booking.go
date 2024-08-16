package model

import "time"

// swagger: model Booking
type Booking struct {
	ID            int        `json:"id" gorm:"primaryKey"`
	UserID        int        `json:"user_id"`
	RoomID        int        `json:"room_id"`
	CheckInDate   *time.Time `json:"check_in_date" gorm:"default:NULL"`
	CheckOutDate  *time.Time `json:"check_out_date" gorm:"default:NULL"`
	PaymentMethod string     `json:"payment_method" gorm:"type:varchar(6)"`
	TotalPrice    float64    `json:"total_price"`
	Status        string     `json:"status" gorm:"type:varchar(3)"`
	Note          string     `json:"note" gorm:"type:text"`
	Quantity      int        `json:"quantity" gorm:"type:integer,default:1"`
	Reason        string     `json:"reason" gorm:"type:text,default:NULL"`

	User *User `gorm:"foreignKey:UserID"`
	Room *Room `json:"room" gorm:"foreignKey:RoomID"`

	Base
}

var (
	BookingStatusPending   = "PEN"
	BookingStatusApproved  = "APP"
	BookingStatusRejected  = "REJ"
	BookingStatusCancel    = "CAN"
	BookingStatusInReview  = "REV"
	BookingStatusCompleted = "COM"

	BookingPaymentMethodCash   = "CASH"
	BookingPaymentMethodOnline = "ONLINE"
)

// swagger: model BookingResponse
type BookingResponse struct {
	ID            int           `json:"id"`
	UserID        int           `json:"user_id"`
	User          *UserResponse `json:"user"`
	RoomID        int           `json:"room_id"`
	RoomName      string        `json:"room_name"`
	CheckInDate   time.Time     `json:"check_in_date"`
	CheckOutDate  time.Time     `json:"check_out_date"`
	PaymentMethod string        `json:"payment_method"`
	TotalPrice    float64       `json:"total_price"`
	Status        string        `json:"status"`
	Note          string        `json:"note"`
	Quantity      int           `json:"quantity"`
	CreatedAt     time.Time     `json:"created_at"`
	Reason        string        `json:"reason"`
}

func (b *Booking) ToResponse() *BookingResponse {
	resp := &BookingResponse{
		ID:            b.ID,
		UserID:        b.UserID,
		RoomID:        b.RoomID,
		CheckInDate:   *b.CheckInDate,
		CheckOutDate:  *b.CheckOutDate,
		PaymentMethod: b.PaymentMethod,
		TotalPrice:    b.TotalPrice,
		Status:        b.Status,
		Note:          b.Note,
		Quantity:      b.Quantity,
		CreatedAt:     b.CreatedAt,
		Reason:        b.Reason,
	}

	if b.User != nil {
		resp.User = b.User.ToResponse()
	}

	if b.Room != nil {
		resp.RoomName = b.Room.RoomName
	}

	return resp
}
