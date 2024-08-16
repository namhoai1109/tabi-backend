package model

import (
	"time"

	"github.com/lib/pq"
)

// swagger:model Branch
type Branch struct {
	ID                    int           `json:"id" gorm:"primaryKey"`
	BranchName            string        `json:"branch_name" gorm:"type:varchar(255)"`
	Address               string        `json:"address" gorm:"type:varchar(128)"`
	FullAddress           string        `json:"full_address" gorm:"type:text"`
	ProvinceCity          string        `json:"province_city" gorm:"type:varchar(64)"`
	District              string        `json:"district" gorm:"type:varchar(64)"`
	Ward                  string        `json:"ward" gorm:"type:varchar(64)"`
	Latitude              string        `json:"latitude" gorm:"type:varchar(64)"`
	Longitude             string        `json:"longitude" gorm:"type:varchar(64)"`
	Description           string        `json:"description" gorm:"type:text"`
	ReceptionArea         bool          `json:"reception_area" gorm:"type:boolean"`
	MainFacilities        pq.Int64Array `json:"main_facilities" gorm:"type:integer[]"`
	IsActive              bool          `json:"is_active" gorm:"type:boolean"`
	CompanyID             int           `json:"company_id" gorm:"default:NULL"`
	BranchManagerID       int           `json:"branch_manager_id" gorm:"default:NULL"`
	TypeID                int           `json:"type_id" gorm:"default:NULL"`
	CancellationTimeUnit  string        `json:"cancellation_time_unit"`
	CancellationTimeValue int           `json:"cancellation_time_value"`
	GeneralPolicy         string        `json:"general_policy"`
	WebsiteURL            string        `json:"website_url" gorm:"default:NULL"`
	TaxNumber             string        `json:"tax_number" gorm:"default:NULL"`

	Company       *Company            `gorm:"foreignKey:CompanyID"`
	BranchManager *BranchManager      `gorm:"foreignKey:BranchManagerID"`
	Type          *GeneralType        `gorm:"foreignKey:TypeID"`
	Banks         []*Bank             `gorm:"foreignKey:BranchID"`
	RoomTypes     []*RoomTypeOfBranch `gorm:"foreignKey:BranchID"`
	Rooms         []*Room             `gorm:"foreignKey:BranchID"`
	Ratings       []*Rating           `json:"ratings" gorm:"foreignKey:BranchID"`
	Base
}

var (
	BranchCancellationTimeUnitHour  = "HOUR"
	BranchCancellationTimeUnitDay   = "DAY"
	BranchCancellationTimeUnitWeek  = "WEEK"
	BranchCancellationTimeUnitMonth = "MONTH"
	BranchCancellationTimeUnitYear  = "YEAR"
)

var BranchCancellationTimeUnits = []string{
	BranchCancellationTimeUnitHour,
	BranchCancellationTimeUnitDay,
	BranchCancellationTimeUnitWeek,
	BranchCancellationTimeUnitMonth,
	BranchCancellationTimeUnitYear,
}

// swagger:model BranchResponse
type BranchResponse struct {
	ID                    int         `json:"id"`
	BranchName            string      `json:"branch_name"`
	Address               string      `json:"address"`
	FullAddress           string      `json:"full_address"`
	ProvinceCity          string      `json:"province_city"`
	District              string      `json:"district"`
	Ward                  string      `json:"ward"`
	Latitude              string      `json:"latitude"`
	Longitude             string      `json:"longitude"`
	Description           string      `json:"description"`
	ReceptionArea         bool        `json:"reception_area"`
	IsActive              bool        `json:"is_active"`
	MainFacilities        []*Facility `json:"main_facilities"`
	CancellationTimeUnit  string      `json:"cancellation_time_unit"`
	CancellationTimeValue int         `json:"cancellation_time_value"`
	GeneralPolicy         string      `json:"general_policy"`
	WebsiteURL            string      `json:"website_url,omitempty"`
	TaxNumber             string      `json:"tax_number,omitempty"`

	BranchManager *BranchManagerResponse `json:"branch_manager"`
	RoomTypes     []*RoomType            `json:"room_types"`
	Type          *GeneralType           `json:"type"`
	Banks         []*Bank                `json:"banks"`
	Ratings       []*RatingResponse      `json:"ratings"`
}

// swagger:model PublicBranchResponse
type PublicBranchResponse struct {
	ID                    int               `json:"id"`
	BranchName            string            `json:"branch_name"`
	Address               string            `json:"address"`
	FullAddress           string            `json:"full_address"`
	ProvinceCity          string            `json:"province_city"`
	District              string            `json:"district"`
	Ward                  string            `json:"ward"`
	Latitude              string            `json:"latitude"`
	Longitude             string            `json:"longitude"`
	Description           string            `json:"description"`
	ReceptionArea         bool              `json:"reception_area"`
	MainFacilities        []*Facility       `json:"main_facilities"`
	Type                  *GeneralType      `json:"type"`
	CancellationTimeUnit  string            `json:"cancellation_time_unit"`
	CancellationTimeValue int               `json:"cancellation_time_value"`
	GeneralPolicy         string            `json:"general_policy"`
	Ratings               []*RatingResponse `json:"ratings"`
	HasPaypal             bool              `json:"has_paypal"`
}

// swagger:model RevenueAnalysisData
type RevenueAnalysisData struct {
	Month   int     `json:"month"`
	Revenue float64 `json:"revenue"`
}

//swagger:model BookingRequestQuantityAnalysisData
type BookingRequestQuantityAnalysisData struct {
	Month    int `json:"month"`
	Quantity int `json:"quantity"`
}

func (s *Branch) ToBranchResponse(facilities []*Facility) *BranchResponse {
	branchManagerResp := &BranchManagerResponse{}
	if s.BranchManager != nil {
		branchManagerResp = s.BranchManager.ToResponse()
	}

	ratingResponses := []*RatingResponse{}

	if s.Ratings == nil {
		s.Ratings = []*Rating{}
	}

	for _, rating := range s.Ratings {
		ratingResponses = append(ratingResponses, rating.ToRatingResponse())
	}

	return &BranchResponse{
		ID:                    s.ID,
		BranchName:            s.BranchName,
		Address:               s.Address,
		FullAddress:           s.FullAddress,
		ProvinceCity:          s.ProvinceCity,
		District:              s.District,
		Ward:                  s.Ward,
		Latitude:              s.Latitude,
		Longitude:             s.Longitude,
		Description:           s.Description,
		ReceptionArea:         s.ReceptionArea,
		MainFacilities:        facilities,
		IsActive:              s.IsActive,
		BranchManager:         branchManagerResp,
		Type:                  s.Type,
		Banks:                 s.Banks,
		CancellationTimeUnit:  s.CancellationTimeUnit,
		CancellationTimeValue: s.CancellationTimeValue,
		GeneralPolicy:         s.GeneralPolicy,
		Ratings:               ratingResponses,
		WebsiteURL:            s.WebsiteURL,
		TaxNumber:             s.TaxNumber,
	}
}

func (s *Branch) ToPublicBranchResponse(facilities []*Facility) *PublicBranchResponse {
	branchType := &GeneralType{}
	if s.Type != nil {
		branchType = s.Type
	}

	ratingResponses := []*RatingResponse{}

	if s.Ratings == nil {
		s.Ratings = []*Rating{}
	}

	for _, rating := range s.Ratings {
		ratingResponses = append(ratingResponses, rating.ToRatingResponse())
	}

	hasPaypal := false
	validHst := s.BranchManager != nil && s.BranchManager.Account != nil
	validRp := s.Company != nil && s.Company.Representative != nil && s.Company.Representative.Account != nil
	if validHst && s.BranchManager.Account.Role == AccountRoleHost && s.BranchManager.Account.Email != "" {
		hasPaypal = true
	} else if validRp && s.Company.Representative.Account.Email != "" {
		hasPaypal = true
	}

	return &PublicBranchResponse{
		ID:                    s.ID,
		BranchName:            s.BranchName,
		Address:               s.Address,
		FullAddress:           s.FullAddress,
		ProvinceCity:          s.ProvinceCity,
		District:              s.District,
		Ward:                  s.Ward,
		Latitude:              s.Latitude,
		Longitude:             s.Longitude,
		Description:           s.Description,
		ReceptionArea:         s.ReceptionArea,
		MainFacilities:        facilities,
		Type:                  branchType,
		CancellationTimeUnit:  s.CancellationTimeUnit,
		CancellationTimeValue: s.CancellationTimeValue,
		GeneralPolicy:         s.GeneralPolicy,
		Ratings:               ratingResponses,
		HasPaypal:             hasPaypal,
	}
}

func (s *Branch) GetCancellationTime() time.Duration {
	switch s.CancellationTimeUnit {
	case BranchCancellationTimeUnitHour:
		return time.Duration(s.CancellationTimeValue) * time.Hour
	case BranchCancellationTimeUnitDay:
		return time.Duration(s.CancellationTimeValue) * 24 * time.Hour
	case BranchCancellationTimeUnitWeek:
		return time.Duration(s.CancellationTimeValue) * 7 * 24 * time.Hour
	case BranchCancellationTimeUnitMonth:
		return time.Duration(s.CancellationTimeValue) * 30 * 24 * time.Hour
	case BranchCancellationTimeUnitYear:
		return time.Duration(s.CancellationTimeValue) * 365 * 24 * time.Hour
	default:
		return 0
	}
}
