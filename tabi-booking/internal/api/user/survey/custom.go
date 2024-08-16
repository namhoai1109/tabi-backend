package survey

//swagger:model SurveyCreationRequest
type SurveyCreationRequest struct {
	PlaceType  string `json:"place_type" validate:"required"`
	Activities string `json:"activities" validate:"required"`
	Seasons    string `json:"seasons" validate:"required"`
}
