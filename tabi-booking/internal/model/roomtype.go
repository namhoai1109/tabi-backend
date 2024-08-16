package model

import (
	"time"

	"github.com/lib/pq"
)

// swagger:model RoomType
type RoomType struct {
	ID               int           `json:"id" gorm:"primaryKey"`
	TypeName         string        `json:"type_name" gorm:"type:varchar(32)"`
	CheckInTime      time.Time     `json:"check_in_time" gorm:"type:time"`
	CheckOutTime     time.Time     `json:"check_out_time" gorm:"type:time"`
	IncludeBreakfast bool          `json:"include_breakfast" gorm:"type:boolean"`
	RoomFacilities   pq.Int64Array `json:"room_facilities" gorm:"type:integer[]"`

	Rooms    []*Room             `gorm:"foreignKey:RoomTypeID"`
	Branches []*RoomTypeOfBranch `gorm:"foreignKey:RoomTypeID"`
	Base
}

// swagger:model RoomTypeResponse
type RoomTypeResponse struct {
	ID               int         `json:"id"`
	TypeName         string      `json:"type_name"`
	CheckInTime      time.Time   `json:"check_in_time"`
	CheckOutTime     time.Time   `json:"check_out_time"`
	IncludeBreakfast bool        `json:"include_breakfast"`
	Facilities       []*Facility `json:"facilities"`
}

func (s *RoomType) ToResponse(facilities []*Facility) *RoomTypeResponse {
	return &RoomTypeResponse{
		ID:               s.ID,
		TypeName:         s.TypeName,
		CheckInTime:      s.CheckInTime,
		CheckOutTime:     s.CheckOutTime,
		IncludeBreakfast: s.IncludeBreakfast,
		Facilities:       facilities,
	}
}
