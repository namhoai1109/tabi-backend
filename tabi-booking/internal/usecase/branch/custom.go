package branch

import (
	"time"

	"github.com/lib/pq"
)

// swagger:model PublicBranch
type PublicBranch struct {
	ID             int     `json:"id"`
	Name           string  `json:"name"`
	District       string  `json:"district"`
	ProvinceCity   string  `json:"province_city"`
	MinPrice       float64 `json:"min_price"`
	MaxPrice       float64 `json:"max_price"`
	StarLevel      float64 `json:"star_level,omitempty"`
	ReviewQuantity int     `json:"review_quantity,omitempty"`
}

type PublicBranchFilter struct {
	BookingDateIn string `json:"booking_date__in,omitempty" query:"booking_date__in"`
	Occupancy     int    `json:"occupancy,omitempty" query:"occupancy"`
}

type PublicBranchCondition struct {
	BookingDateIn []time.Time `json:"booking_date__in,omitempty"`
	Occupancy     *int        `json:"occupancy,omitempty"`
}

type PublicBranchListResponse struct {
	Data  []*PublicBranch `json:"data"`
	Total int64           `json:"total"`
}

// BranchCreationRequest represents branch creation request
// swagger:model BranchCreationRequest
type BranchCreationRequest struct {
	BranchName            string        `json:"branch_name" validate:"required"`
	Address               string        `json:"address" validate:"required"`
	ProvinceCity          string        `json:"province_city" validate:"required"`
	District              string        `json:"district" validate:"required"`
	Ward                  string        `json:"ward" validate:"required"`
	Latitude              string        `json:"latitude,omitempty"`
	Longitude             string        `json:"longitude,omitempty"`
	Description           string        `json:"description,omitempty"`
	ReceptionArea         bool          `json:"reception_area"`
	MainFacilities        pq.Int64Array `json:"main_facilities" validate:"required"`
	TypeID                int           `json:"type_id" validate:"required"`
	CancellationTimeUnit  string        `json:"cancellation_time_unit" validate:"required,min=1"`
	CancellationTimeValue int           `json:"cancellation_time_value" validate:"required"`
	GeneralPolicy         string        `json:"general_policy" validate:"required"`
	BranchManagerID       *int          `json:"branch_manager_id"`
	WebsiteURL            string        `json:"website_url,omitempty"`
	TaxNumber             string        `json:"tax_number,omitempty"`
}
