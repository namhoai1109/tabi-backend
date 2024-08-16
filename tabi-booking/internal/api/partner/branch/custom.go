package branch

import "github.com/lib/pq"

type TypeResponse struct {
	ID      int    `json:"id"`
	LabelVI string `json:"label_vi"`
	LabelEN string `json:"label_en"`
	DescVI  string `json:"desc_vi"`
	DescEN  string `json:"desc_en"`
	Order   int    `json:"order"`
}

type CustomBranchResponse struct {
	ID            int           `json:"id"`
	BranchName    string        `json:"branch_name"`
	Address       string        `json:"address"`
	FullAddress   string        `json:"full_address"`
	ProvinceCity  string        `json:"province_city"`
	District      string        `json:"district"`
	Ward          string        `json:"ward"`
	ReceptionArea bool          `json:"reception_area"`
	Description   string        `json:"description"`
	TypeResponse  *TypeResponse `json:"type"`
}

// swagger:model BranchListResponse
type BranchListResponse struct {
	Total int64                   `json:"total"`
	Data  []*CustomBranchResponse `json:"data"`
}

// swagger:model BranchUpdateRequest
type BranchUpdateRequest struct {
	BranchName            string        `json:"branch_name,omitempty"`
	Address               string        `json:"address,omitempty"`
	ProvinceCity          string        `json:"province_city,omitempty"`
	District              string        `json:"district,omitempty"`
	Ward                  string        `json:"ward,omitempty"`
	Latitude              string        `json:"latitude,omitempty"`
	Longitude             string        `json:"longitude,omitempty"`
	Description           string        `json:"description,omitempty"`
	ReceptionArea         bool          `json:"reception_area,omitempty"`
	MainFacilities        pq.Int64Array `json:"main_facilities,omitempty"`
	TypeID                int           `json:"type_id,omitempty"`
	CancellationTimeUnit  string        `json:"cancellation_time_unit,omitempty"`
	CancellationTimeValue int           `json:"cancellation_time_value,omitempty"`
	GeneralPolicy         string        `json:"general_policy,omitempty"`
	WebsiteURL            string        `json:"website_url,omitempty"`
	TaxNumber             string        `json:"tax_number,omitempty"`
	Email                 *string       `json:"email,omitempty"`
}

const (
	LimitPerPage = 10
)

// swagger:model ActivateBranchResponse
type ActivateBranchResponse struct {
	Activated bool `json:"activated"`
}

// swagger:parameters PartnerBranchAnalysisRevenues PartnerBranchAnalysisBookingRequestQuantity
type BranchAnalysisRequest struct {
	Year int `json:"year,omitempty" query:"year"`
}
