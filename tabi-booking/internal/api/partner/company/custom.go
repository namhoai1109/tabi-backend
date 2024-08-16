package company

// swagger:model CompanyUpdateRequest
type CompanyUpdateRequest struct {
	CompanyName *string `json:"company_name,omitempty"`
	ShortName   *string `json:"short_name,omitempty"`
	Description *string `json:"description,omitempty"`
	WebsiteURL  *string `json:"website_url,omitempty"`
	TaxNumber   *string `json:"tax_number,omitempty"`
	Email       *string `json:"email,omitempty"`
}

// swagger:parameters PartnerCompanyAnalysisRevenues PartnerCompanyAnalysisBookingRequestQuantity
type CompanyAnalysisRequest struct {
	Year int `json:"year,omitempty" query:"year"`
}
