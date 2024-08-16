package model

type ReservationReduction struct {
	ID        int     `json:"id" gorm:"primaryKey"`
	RoomID    int     `json:"room_id"`
	Quantity  float64 `json:"quantity"`
	TimeUnit  string  `json:"time_unit" gorm:"type:varchar(5)"`
	Reduction float64 `json:"reduction"`

	Room *Room `gorm:"foreignKey:RoomID"`
	Base
}

var (
	ReservationReductionTimeUnitHour  = "HOUR"
	ReservationReductionTimeUnitDay   = "DAY"
	ReservationReductionTimeUnitWeek  = "WEEK"
	ReservationReductionTimeUnitMonth = "MONTH"
	ReservationReductionTimeUnitYear  = "YEAR"
)

var ReservationReductionTimeUnits = []string{
	ReservationReductionTimeUnitHour,
	ReservationReductionTimeUnitDay,
	ReservationReductionTimeUnitWeek,
	ReservationReductionTimeUnitMonth,
	ReservationReductionTimeUnitYear,
}

func (s *ReservationReduction) convertToDateQuantity() float64 {
	switch s.TimeUnit {
	case ReservationReductionTimeUnitHour:
		return s.Quantity / 24
	case ReservationReductionTimeUnitDay:
		return s.Quantity
	case ReservationReductionTimeUnitWeek:
		return s.Quantity * 7
	case ReservationReductionTimeUnitMonth:
		return s.Quantity * 30
	case ReservationReductionTimeUnitYear:
		return s.Quantity * 365
	}

	return 0
}
