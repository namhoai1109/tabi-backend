package generaltype

import httpcore "github.com/namhoai1109/tabi/core/http"

//swagger:model AccommodationTypeListResponse
type AccommodationTypeListResponse struct {
	Data []*AccommodationTypeResponse `json:"data"`
}

//swagger:model AccommodationTypeResponse
type AccommodationTypeResponse struct {
	ID          int                          `json:"id"`
	Label       string                       `json:"label"`
	Description string                       `json:"description"`
	Children    []*AccommodationTypeChildren `json:"children"`
}

//swagger:model AccommodationTypeChildren
type AccommodationTypeChildren struct {
	ID          int    `json:"id"`
	Label       string `json:"label"`
	Description string `json:"description"`
}

// swagger:model BedTypeResponse
type BedTypeResponse struct {
	ID    int    `json:"id"`
	Label string `json:"label"`
}

// swagger:model BedTypeListResponse
type BedTypeListResponse struct {
	Data []*BedTypeResponse `json:"data"`
}

// swagger:parameters PublicGeneralTypeAccommodationList PublicGeneralTypeBedList
type ListRequestCustom struct {
	httpcore.ListRequest
	// language vi or en
	// default: en
	// required: true
	// in: path
	Lang string `json:"lang" query:"lang"`
}
