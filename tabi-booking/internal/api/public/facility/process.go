package facility

import (
	"tabi-booking/internal/model"
	"tabi-booking/internal/util"

	"golang.org/x/exp/slices"
)

func (s *Facility) transformToFacilityResponse(facilities []*model.Facility, lang string) []*FacilityResponse {

	classes := []string{}
	for _, facility := range facilities {
		class := util.TernaryOperator(lang == model.FacilityLangVI, facility.ClassVI, facility.ClassEN).(string)
		if !slices.Contains(classes, class) {
			classes = append(classes, class)
		}
	}

	facilityGroup := map[string][]*Items{}
	for _, facility := range facilities {

		class := util.TernaryOperator(lang == model.FacilityLangVI, facility.ClassVI, facility.ClassEN).(string)
		name := util.TernaryOperator(lang == model.FacilityLangVI, facility.NameVI, facility.NameEN).(string)

		facilityGroup[class] = append(facilityGroup[class], &Items{
			ID:   facility.ID,
			Name: name,
		})
	}

	facilitiesResp := []*FacilityResponse{}
	for _, class := range classes {
		facilitiesResp = append(facilitiesResp, &FacilityResponse{
			Class: class,
			Items: facilityGroup[class],
		})
	}

	return facilitiesResp
}
