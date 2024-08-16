package model

import (
	"time"
)

// swagger:model Room
type Room struct {
	ID           int     `json:"id" gorm:"primaryKey"`
	Status       string  `json:"status" gorm:"type:varchar(3)"`
	RoomName     string  `json:"room_name" gorm:"type:varchar(32)"`
	MaxOccupancy int     `json:"max_occupancy" gorm:"type:integer"`
	Width        float32 `json:"width" gorm:"type:float"`
	Length       float32 `json:"length" gorm:"type:float"`
	MaxPrice     float64 `json:"max_price" gorm:"type:float"`
	Quantity     int     `json:"quantity" gorm:"type:integer"`
	RoomTypeID   int     `json:"room_type_id"`
	BedTypeID    int     `json:"bed_type_id"`
	BranchID     int     `json:"branch_id"`

	RoomType             *RoomType               `json:"room_type" gorm:"foreignKey:RoomTypeID"`
	BedType              *GeneralType            `json:"bed_type" gorm:"foreignKey:BedTypeID"`
	FactureReduction     *FactureReduction       `json:"facture_reduction" gorm:"foreignKey:RoomID"`
	ReservationReduction []*ReservationReduction `json:"reservation_reduction" gorm:"foreignKey:RoomID"`
	Branch               *Branch                 `json:"branch" gorm:"foreignKey:BranchID"`
	Bookings             []*Booking              `json:"bookings" gorm:"foreignKey:RoomID"`
	Base
}

var (
	RoomStatusPending  = "PEN" // Pending
	RoomStatusApproved = "APP" // Approved (Approved but has no room images)
	RoomStatusRejected = "REJ" // Rejected
	RoomStatusUpdated  = "UPD" // Updated (Added room images)
)

// swagger:model PublicRoom
type PublicRoom struct {
	ID           int     `json:"id"`
	RoomName     string  `json:"room_name"`
	MaxOccupancy int     `json:"max_occupancy"`
	Width        float32 `json:"width"`
	Length       float32 `json:"length"`

	MaxPrice              float64 `json:"max_price"`
	CurrentPrice          float64 `json:"current_price"`
	OnlineMethodReduction float64 `json:"online_method_reduction"`
	OnCashMethodReduction float64 `json:"on_cash_method_reduction"`
	RemainingQuantity     int     `json:"remaining_quantity"`

	RoomType *RoomTypeResponse `json:"room_type" gorm:"foreignKey:RoomTypeID"`
	BedType  *GeneralType      `json:"bed_type" gorm:"foreignKey:BedTypeID"`
}

func (s *Room) ToPublicRoomResponse(facilities []*Facility, checkInDate, checkOutDate time.Time) *PublicRoom {
	remainingQuantity := s.CountAvailableRoom(checkInDate, checkOutDate)
	if remainingQuantity < 0 {
		remainingQuantity = 0
	}
	resp := &PublicRoom{
		ID:                s.ID,
		RoomName:          s.RoomName,
		MaxOccupancy:      s.MaxOccupancy,
		Width:             s.Width,
		Length:            s.Length,
		MaxPrice:          s.MaxPrice,
		RemainingQuantity: remainingQuantity,
	}

	if s.RoomType != nil {
		resp.RoomType = s.RoomType.ToResponse(facilities)
	}

	if s.BedType != nil {
		resp.BedType = s.BedType
	}

	if s.FactureReduction != nil && s.ReservationReduction != nil {
		resp.OnlineMethodReduction = s.FactureReduction.OnlineMethod
		resp.OnCashMethodReduction = s.FactureReduction.OnCashMethod
		resp.CurrentPrice = s.GetPriceForBookingDates(checkInDate, checkOutDate)
	}

	return resp
}

func (s *Room) CountAvailableRoom(checkInDate, checkOutDate time.Time) int {
	if s.Bookings == nil || len(s.Bookings) == 0 {
		return s.Quantity
	}

	count := 0
	for _, booking := range s.Bookings {
		if s.ID == booking.RoomID && (booking.Status == BookingStatusApproved || booking.Status == BookingStatusPending) {
			if booking.CheckInDate.Before(checkInDate) && booking.CheckOutDate.After(checkInDate) {
				count += booking.Quantity
				continue
			}

			if booking.CheckInDate.Before(checkOutDate) && booking.CheckOutDate.After(checkOutDate) {
				count += booking.Quantity
				continue
			}

			if booking.CheckInDate.After(checkInDate) && booking.CheckOutDate.Before(checkOutDate) {
				count += booking.Quantity
				continue
			}

			if booking.CheckInDate.Equal(checkInDate) || booking.CheckOutDate.Equal(checkOutDate) {
				count += booking.Quantity
				continue
			}
		}
	}

	return s.Quantity - count
}

func (s *Room) getDates(start, end time.Time) []time.Time {
	var dates []time.Time
	for d := start; d.Before(end); d = d.AddDate(0, 0, 1) {
		dates = append(dates, d)
	}
	return dates
}

func (s *Room) findReservationReduction(date time.Time) float64 {
	reservationTime := time.Until(date).Hours() / 24
	reservationReduction := float64(0)
	for _, reduction := range s.ReservationReduction {
		if reservationTime >= reduction.convertToDateQuantity() {
			reservationReduction = reduction.Reduction
		}
	}

	return reservationReduction
}

func (s *Room) GetPriceForBookingDates(checkInDate, checkOutDate time.Time) float64 {
	averagePrice := float64(0)
	if s.FactureReduction == nil && s.ReservationReduction == nil {
		return 0
	}

	bookingDates := s.getDates(checkInDate, checkOutDate)
	for _, date := range bookingDates {
		reservationReduction := s.findReservationReduction(date)
		factureReduction := s.FactureReduction.GetReduction(&date)

		averagePrice += s.MaxPrice * (1 - (reservationReduction + factureReduction))
	}
	averagePrice = averagePrice / float64(len(bookingDates))

	return averagePrice
}
