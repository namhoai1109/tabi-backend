package model

import (
	"strings"
	"time"
)

type Holiday struct {
	NameEN string `json:"name_en"`
	NameVI string `json:"name_vi"`
	Start  string `json:"start"`
	End    string `json:"end"`
}

var holidays = []Holiday{
	{
		NameEN: "International New Year's Day",
		NameVI: "Ngày Quốc tế Năm mới",
		Start:  "2024-01-01",
		End:    "2024-01-02",
	},

	{
		NameEN: "Vietnamese New Year's Eve",
		NameVI: "Đêm giao thừa của người Việt",
		Start:  "2024-02-09",
		End:    "2024-02-10",
	},

	{
		NameEN: "Vietnamese New Year",
		NameVI: "Tết Việt Nam",
		Start:  "2024-02-10",
		End:    "2024-02-11",
	},

	{
		NameEN: "Tet holiday",
		NameVI: "Tết",
		Start:  "2024-02-11",
		End:    "2024-02-12",
	},

	{
		NameEN: "Tet holiday",
		NameVI: "Tết",
		Start:  "2024-02-12",
		End:    "2024-02-13",
	},

	{
		NameEN: "Tet holiday",
		NameVI: "Tết",
		Start:  "2024-02-13",
		End:    "2024-02-14",
	},

	{
		NameEN: "Tet holiday",
		NameVI: "Tết",
		Start:  "2024-02-14",
		End:    "2024-02-15",
	},

	{
		NameEN: "Easter Sunday",
		NameVI: "Chủ Nhật Phục Sinh",
		Start:  "2024-03-31",
		End:    "2024-04-01",
	},

	{
		NameEN: "Hung Kings Festival",
		NameVI: "Lễ hội các Vua Hùng",
		Start:  "2024-04-18",
		End:    "2024-04-19",
	},

	{
		NameEN: "Liberation Day/Reunification Day",
		NameVI: "Ngày giải phóng/Ngày thống nhất",
		Start:  "2024-04-30",
		End:    "2024-05-01",
	},

	{
		NameEN: "International Labor Day",
		NameVI: "Ngày quốc tế lao động",
		Start:  "2024-05-01",
		End:    "2024-05-02",
	},

	{
		NameEN: "Vesak",
		NameVI: "Lễ Vesak",
		Start:  "2024-05-23",
		End:    "2024-05-24",
	},

	{
		NameEN: "Independence Day",
		NameVI: "Ngày Quốc Khánh",
		Start:  "2024-09-02",
		End:    "2024-09-03",
	},

	{
		NameEN: "Christmas Eve",
		NameVI: "Đêm Giáng sinh",
		Start:  "2024-12-24",
		End:    "2024-12-25",
	},

	{
		NameEN: "Christmas Day",
		NameVI: "Ngày Giáng Sinh",
		Start:  "2024-12-25",
		End:    "2024-12-26",
	},

	{
		NameEN: "International New Year's Eve",
		NameVI: "Đêm giao thừa quốc tế",
		Start:  "2024-12-31",
		End:    "2025-01-01",
	},

	{
		NameEN: "International New Year's Day",
		NameVI: "Ngày Quốc tế Năm mới",
		Start:  "2025-01-01",
		End:    "2025-01-02",
	},

	{
		NameEN: "Vietnamese New Year's Eve",
		NameVI: "Đêm giao thừa của người Việt",
		Start:  "2025-01-28",
		End:    "2025-01-29",
	},

	{
		NameEN: "Vietnamese New Year",
		NameVI: "Tết Việt Nam",
		Start:  "2025-01-29",
		End:    "2025-01-30",
	},

	{
		NameEN: "Tet holiday",
		NameVI: "Tết",
		Start:  "2025-01-30",
		End:    "2025-01-31",
	},

	{
		NameEN: "Tet holiday",
		NameVI: "Tết",
		Start:  "2025-01-31",
		End:    "2025-02-01",
	},

	{
		NameEN: "Tet holiday",
		NameVI: "Tết",
		Start:  "2025-02-01",
		End:    "2025-02-02",
	},

	{
		NameEN: "Tet holiday",
		NameVI: "Tết",
		Start:  "2025-02-02",
		End:    "2025-02-03",
	},

	{
		NameEN: "Hung Kings Festival",
		NameVI: "Lễ hội các Vua Hùng",
		Start:  "2025-04-07",
		End:    "2025-04-08",
	},

	{
		NameEN: "Easter Sunday",
		NameVI: "Chủ Nhật Phục Sinh",
		Start:  "2025-04-20",
		End:    "2025-04-21",
	},

	{
		NameEN: "Liberation Day/Reunification Day",
		NameVI: "Ngày giải phóng/Ngày thống nhất",
		Start:  "2025-04-30",
		End:    "2025-05-01",
	},

	{
		NameEN: "International Labor Day",
		NameVI: "Ngày quốc tế lao động",
		Start:  "2025-05-01",
		End:    "2025-05-02",
	},

	{
		NameEN: "Vesak",
		NameVI: "Lễ Vesak",
		Start:  "2025-05-12",
		End:    "2025-05-13",
	},

	{
		NameEN: "Independence Day",
		NameVI: "Ngày Quốc Khánh",
		Start:  "2025-09-02",
		End:    "2025-09-03",
	},

	{
		NameEN: "Christmas Eve",
		NameVI: "Đêm Giáng sinh",
		Start:  "2025-12-24",
		End:    "2025-12-25",
	},

	{
		NameEN: "Christmas Day",
		NameVI: "Ngày Giáng Sinh",
		Start:  "2025-12-25",
		End:    "2025-12-26",
	},

	{
		NameEN: "International New Year's Eve",
		NameVI: "Đêm giao thừa quốc tế",
		Start:  "2025-12-31",
		End:    "2026-01-01",
	},

	{
		NameEN: "International New Year's Day",
		NameVI: "Ngày Quốc tế Năm mới",
		Start:  "2026-01-01",
		End:    "2026-01-02",
	},

	{
		NameEN: "Vietnamese New Year's Eve",
		NameVI: "Đêm giao thừa của người Việt",
		Start:  "2026-02-16",
		End:    "2026-02-17",
	},

	{
		NameEN: "Vietnamese New Year",
		NameVI: "Tết Việt Nam",
		Start:  "2026-02-17",
		End:    "2026-02-18",
	},

	{
		NameEN: "Tet holiday",
		NameVI: "Tết",
		Start:  "2026-02-18",
		End:    "2026-02-19",
	},

	{
		NameEN: "Tet holiday",
		NameVI: "Tết",
		Start:  "2026-02-19",
		End:    "2026-02-20",
	},

	{
		NameEN: "Tet holiday",
		NameVI: "Tết",
		Start:  "2026-02-20",
		End:    "2026-02-21",
	},

	{
		NameEN: "Tet holiday",
		NameVI: "Tết",
		Start:  "2026-02-21",
		End:    "2026-02-22",
	},

	{
		NameEN: "Easter Sunday",
		NameVI: "Chủ Nhật Phục Sinh",
		Start:  "2026-04-05",
		End:    "2026-04-06",
	},

	{
		NameEN: "Hung Kings Festival",
		NameVI: "Lễ hội các Vua Hùng",
		Start:  "2026-04-26",
		End:    "2026-04-27",
	},

	{
		NameEN: "Liberation Day/Reunification Day",
		NameVI: "Ngày giải phóng/Ngày thống nhất",
		Start:  "2026-04-30",
		End:    "2026-05-01",
	},

	{
		NameEN: "International Labor Day",
		NameVI: "Ngày quốc tế lao động",
		Start:  "2026-05-01",
		End:    "2026-05-02",
	},

	{
		NameEN: "Independence Day",
		NameVI: "Ngày Quốc Khánh",
		Start:  "2026-09-02",
		End:    "2026-09-03",
	},

	{
		NameEN: "Christmas Eve",
		NameVI: "Đêm Giáng sinh",
		Start:  "2026-12-24",
		End:    "2026-12-25",
	},

	{
		NameEN: "Christmas Day",
		NameVI: "Ngày Giáng Sinh",
		Start:  "2026-12-25",
		End:    "2026-12-26",
	},

	{
		NameEN: "International New Year's Eve",
		NameVI: "Đêm giao thừa quốc tế",
		Start:  "2026-12-31",
		End:    "2027-01-01",
	},

	{
		NameEN: "International New Year's Day",
		NameVI: "Ngày Quốc tế Năm mới",
		Start:  "2027-01-01",
		End:    "2027-01-02",
	},

	{
		NameEN: "Vietnamese New Year's Eve",
		NameVI: "Đêm giao thừa của người Việt",
		Start:  "2027-02-05",
		End:    "2027-02-06",
	},

	{
		NameEN: "Vietnamese New Year",
		NameVI: "Tết Việt Nam",
		Start:  "2027-02-06",
		End:    "2027-02-07",
	},

	{
		NameEN: "Tet holiday",
		NameVI: "Tết",
		Start:  "2027-02-07",
		End:    "2027-02-08",
	},

	{
		NameEN: "Tet holiday",
		NameVI: "Tết",
		Start:  "2027-02-08",
		End:    "2027-02-09",
	},

	{
		NameEN: "Tet holiday",
		NameVI: "Tết",
		Start:  "2027-02-09",
		End:    "2027-02-10",
	},

	{
		NameEN: "Tet holiday",
		NameVI: "Tết",
		Start:  "2027-02-10",
		End:    "2027-02-11",
	},

	{
		NameEN: "Easter Sunday",
		NameVI: "Chủ Nhật Phục Sinh",
		Start:  "2027-03-28",
		End:    "2027-03-29",
	},

	{
		NameEN: "Hung Kings Festival",
		NameVI: "Lễ hội các Vua Hùng",
		Start:  "2027-04-16",
		End:    "2027-04-17",
	},

	{
		NameEN: "Liberation Day/Reunification Day",
		NameVI: "Ngày giải phóng/Ngày thống nhất",
		Start:  "2027-04-30",
		End:    "2027-05-01",
	},

	{
		NameEN: "International Labor Day",
		NameVI: "Ngày quốc tế lao động",
		Start:  "2027-05-01",
		End:    "2027-05-02",
	},

	{
		NameEN: "Day off for International Labor Day",
		NameVI: "Ngày nghỉ nhân Ngày Quốc tế Lao động",
		Start:  "2027-05-03",
		End:    "2027-05-04",
	},

	{
		NameEN: "Independence Day",
		NameVI: "Ngày Quốc Khánh",
		Start:  "2027-09-02",
		End:    "2027-09-03",
	},

	{
		NameEN: "Christmas Eve",
		NameVI: "Đêm Giáng sinh",
		Start:  "2027-12-24",
		End:    "2027-12-25",
	},

	{
		NameEN: "Christmas Day",
		NameVI: "Ngày Giáng Sinh",
		Start:  "2027-12-25",
		End:    "2027-12-26",
	},

	{
		NameEN: "International New Year's Eve",
		NameVI: "Đêm giao thừa quốc tế",
		Start:  "2027-12-31",
		End:    "2028-01-01",
	},

	{
		NameEN: "International New Year's Day",
		NameVI: "Ngày Quốc tế Năm mới",
		Start:  "2028-01-01",
		End:    "2028-01-02",
	},

	{
		NameEN: "Day off for International New Year's Day",
		NameVI: "Ngày nghỉ Quốc tế Năm mới",
		Start:  "2028-01-03",
		End:    "2028-01-04",
	},

	{
		NameEN: "Vietnamese New Year's Eve",
		NameVI: "Đêm giao thừa của người Việt",
		Start:  "2028-01-25",
		End:    "2028-01-26",
	},

	{
		NameEN: "Vietnamese New Year",
		NameVI: "Tết Việt Nam",
		Start:  "2028-01-26",
		End:    "2028-01-27",
	},

	{
		NameEN: "Tet holiday",
		NameVI: "Tết",
		Start:  "2028-01-27",
		End:    "2028-01-28",
	},

	{
		NameEN: "Tet holiday",
		NameVI: "Tết",
		Start:  "2028-01-28",
		End:    "2028-01-29",
	},

	{
		NameEN: "Tet holiday",
		NameVI: "Tết",
		Start:  "2028-01-29",
		End:    "2028-01-30",
	},

	{
		NameEN: "Tet holiday",
		NameVI: "Tết",
		Start:  "2028-01-30",
		End:    "2028-01-31",
	},

	{
		NameEN: "Hung Kings Festival",
		NameVI: "Lễ hội các Vua Hùng",
		Start:  "2028-04-04",
		End:    "2028-04-05",
	},

	{
		NameEN: "Easter Sunday",
		NameVI: "Chủ Nhật Phục Sinh",
		Start:  "2028-04-16",
		End:    "2028-04-17",
	},

	{
		NameEN: "Liberation Day/Reunification Day",
		NameVI: "Ngày giải phóng/Ngày thống nhất",
		Start:  "2028-04-30",
		End:    "2028-05-01",
	},

	{
		NameEN: "Day off for Liberation Day/Reunification Day",
		NameVI: "Ngày nghỉ giải phóng/Ngày thống nhất",
		Start:  "2028-05-01",
		End:    "2028-05-02",
	},

	{
		NameEN: "International Labor Day",
		NameVI: "Ngày quốc tế lao động",
		Start:  "2028-05-01",
		End:    "2028-05-02",
	},

	{
		NameEN: "Day off for International Labor Day",
		NameVI: "Ngày nghỉ nhân Ngày Quốc tế Lao động",
		Start:  "2028-05-02",
		End:    "2028-05-03",
	},

	{
		NameEN: "Independence Day",
		NameVI: "Ngày Quốc Khánh",
		Start:  "2028-09-02",
		End:    "2028-09-03",
	},

	{
		NameEN: "Independence Day observed",
		NameVI: "Ngày Độc lập được quan sát",
		Start:  "2028-09-04",
		End:    "2028-09-05",
	},

	{
		NameEN: "Christmas Eve",
		NameVI: "Đêm Giáng sinh",
		Start:  "2028-12-24",
		End:    "2028-12-25",
	},

	{
		NameEN: "Christmas Day",
		NameVI: "Ngày Giáng Sinh",
		Start:  "2028-12-25",
		End:    "2028-12-26",
	},

	{
		NameEN: "International New Year's Eve",
		NameVI: "Đêm giao thừa quốc tế",
		Start:  "2028-12-31",
		End:    "2029-01-01",
	},

	{
		NameEN: "Tet holiday",
		NameVI: "Tết",
		Start:  "2024-02-08",
		End:    "2024-02-09",
	},

	{
		NameEN: "Independence Day Holiday",
		NameVI: "Ngày lễ Độc lập",
		Start:  "2024-09-03",
		End:    "2024-09-04",
	},

	{
		NameEN: "International New Year's Day",
		NameVI: "Ngày Quốc tế Năm mới",
		Start:  "2029-01-01",
		End:    "2029-01-02",
	},

	{
		NameEN: "Vietnamese New Year's Eve",
		NameVI: "Đêm giao thừa của người Việt",
		Start:  "2029-02-12",
		End:    "2029-02-13",
	},

	{
		NameEN: "Vietnamese New Year",
		NameVI: "Tết Việt Nam",
		Start:  "2029-02-13",
		End:    "2029-02-14",
	},

	{
		NameEN: "Tet holiday",
		NameVI: "Tết",
		Start:  "2029-02-14",
		End:    "2029-02-15",
	},

	{
		NameEN: "Tet holiday",
		NameVI: "Tết",
		Start:  "2029-02-15",
		End:    "2029-02-16",
	},

	{
		NameEN: "Tet holiday",
		NameVI: "Tết",
		Start:  "2029-02-16",
		End:    "2029-02-17",
	},

	{
		NameEN: "Tet holiday",
		NameVI: "Tết",
		Start:  "2029-02-17",
		End:    "2029-02-18",
	},

	{
		NameEN: "Easter Sunday",
		NameVI: "Chủ Nhật Phục Sinh",
		Start:  "2029-04-01",
		End:    "2029-04-02",
	},

	{
		NameEN: "Hung Kings Festival",
		NameVI: "Lễ hội các Vua Hùng",
		Start:  "2029-04-23",
		End:    "2029-04-24",
	},

	{
		NameEN: "Liberation Day/Reunification Day",
		NameVI: "Ngày giải phóng/Ngày thống nhất",
		Start:  "2029-04-30",
		End:    "2029-05-01",
	},

	{
		NameEN: "International Labor Day",
		NameVI: "Ngày quốc tế lao động",
		Start:  "2029-05-01",
		End:    "2029-05-02",
	},

	{
		NameEN: "Independence Day",
		NameVI: "Ngày Quốc Khánh",
		Start:  "2029-09-02",
		End:    "2029-09-03",
	},

	{
		NameEN: "Independence Day observed",
		NameVI: "Ngày Độc lập được quan sát",
		Start:  "2029-09-03",
		End:    "2029-09-04",
	},

	{
		NameEN: "Christmas Eve",
		NameVI: "Đêm Giáng sinh",
		Start:  "2029-12-24",
		End:    "2029-12-25",
	},

	{
		NameEN: "Christmas Day",
		NameVI: "Ngày Giáng Sinh",
		Start:  "2029-12-25",
		End:    "2029-12-26",
	},

	{
		NameEN: "International New Year's Eve",
		NameVI: "Đêm giao thừa quốc tế",
		Start:  "2029-12-31",
		End:    "2030-01-01",
	},
}

type FactureReduction struct {
	ID           int     `json:"id" gorm:"primaryKey"`
	RoomID       int     `json:"room_id"`
	OnlineMethod float64 `json:"online_method"`
	OnCashMethod float64 `json:"on_cash_method"`
	NormalDay    float64 `json:"normal_day"`
	Holiday      float64 `json:"holiday"`
	Weekend      float64 `json:"weekend"`

	Room *Room `gorm:"foreignKey:RoomID"`
	Base
}

func (s *FactureReduction) isWeekend(day *time.Time) bool {
	t := time.Now()

	if day != nil {
		t = *day
	}

	switch t.Weekday() {
	case time.Saturday, time.Sunday:
		return true
	default:
		return false
	}
}

func (s *FactureReduction) isHoliday(day *time.Time) bool {
	t := time.Now()

	if day != nil {
		t = *day
	}

	tStr := strings.Split(t.String(), " ")[0]
	for _, h := range holidays {
		if tStr >= h.Start && tStr <= h.End {
			return true
		}
	}

	return false
}

func (s *FactureReduction) GetReduction(day *time.Time) float64 {
	switch {
	case s.isHoliday(day):
		return s.Holiday
	case s.isWeekend(day):
		return s.Weekend
	default:
		return s.NormalDay
	}
}
