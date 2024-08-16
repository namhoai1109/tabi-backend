package model

//swagger:model Company
type Company struct {
	ID               int             `json:"id" gorm:"primaryKey"`
	CompanyName      string          `json:"company_name" gorm:"type:varchar(64)"`
	ShortName        string          `json:"short_name" gorm:"type:varchar(24)"`
	Description      string          `json:"description" gorm:"type:text"`
	Address          string          `json:"address" gorm:"type:varchar(128)"`
	FullAddress      string          `json:"full_address" gorm:"type:varchar(128)"`
	ProvinceCity     string          `json:"province_city" gorm:"type:varchar(64)"`
	District         string          `json:"district" gorm:"type:varchar(64)"`
	Ward             string          `json:"ward" gorm:"type:varchar(64)"`
	Latitude         string          `json:"latitude" gorm:"type:varchar(64)"`
	Longitude        string          `json:"longitude" gorm:"type:varchar(64)"`
	WebsiteURL       string          `json:"website_url" gorm:"type:varchar(64)"`
	TaxNumber        string          `json:"tax_number" gorm:"type:varchar(20)"`
	RepresentativeID int             `json:"representative_id"`
	Representative   *Representative `gorm:"foreignKey:RepresentativeID"`
	Branches         []*Branch       `gorm:"foreignKey:CompanyID"`
	Base
}

// swagger:model CompanyResponse
type CompanyResponse struct {
	ID           int                     `json:"id"`
	CompanyName  string                  `json:"company_name"`
	ShortName    string                  `json:"short_name"`
	Description  string                  `json:"description"`
	Address      string                  `json:"address"`
	FullAddress  string                  `json:"full_address"`
	ProvinceCity string                  `json:"province_city"`
	District     string                  `json:"district"`
	Ward         string                  `json:"ward"`
	Latitude     string                  `json:"latitude"`
	Longitude    string                  `json:"longitude"`
	WebsiteURL   string                  `json:"website_url"`
	TaxNumber    string                  `json:"tax_number"`
	RPResponse   *RepresentativeResponse `json:"representative"`
}

func (c *Company) ToCompanyResponse() *CompanyResponse {
	if c.Representative == nil {
		c.Representative = &Representative{}
	}

	return &CompanyResponse{
		ID:           c.ID,
		CompanyName:  c.CompanyName,
		ShortName:    c.ShortName,
		Description:  c.Description,
		Address:      c.Address,
		FullAddress:  c.FullAddress,
		ProvinceCity: c.ProvinceCity,
		District:     c.District,
		Ward:         c.Ward,
		Latitude:     c.Latitude,
		Longitude:    c.Longitude,
		WebsiteURL:   c.WebsiteURL,
		TaxNumber:    c.TaxNumber,
		RPResponse:   c.Representative.ToRepresentativeResponse(),
	}
}
