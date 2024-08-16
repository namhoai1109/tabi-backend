package facility

import httpcore "github.com/namhoai1109/tabi/core/http"

//swagger:model FacilityResponse
type FacilityResponse struct {
	Class string   `json:"class"`
	Items []*Items `json:"items"`
}

//swagger:model FacilityListResponse
type FacilityListResponse struct {
	Data []*FacilityResponse `json:"data"`
}

// swagger:model Items
type Items struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// swagger:parameters PublicFacilityList
type ListRequestCustom struct {
	httpcore.ListRequest
	// language vi or en
	// default: en
	// required: true
	// in: path
	Lang string `json:"lang" query:"lang"`
}
