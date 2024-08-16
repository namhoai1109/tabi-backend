package generaltype

import (
	"context"
	"tabi-booking/internal/model"
	"tabi-booking/internal/util"
)

func (s *GeneralType) ListAccommodation(ctx context.Context, lang string) ([]*AccommodationTypeResponse, error) {

	typeParents, err := s.getTypeByOrder(ctx, model.GeneralTypeAccommodationClass, lang, 1)
	if err != nil {
		return nil, err
	}

	typeChilds, err := s.getTypeByOrder(ctx, model.GeneralTypeAccommodationClass, lang, 2)
	if err != nil {
		return nil, err
	}

	resp := []*AccommodationTypeResponse{}
	isVI := lang == model.GeneralTypeLangVI
	for _, typeParent := range typeParents {
		typeResp := &AccommodationTypeResponse{
			ID:          typeParent.ID,
			Label:       util.TernaryOperator(isVI, typeParent.LabelVI, typeParent.LabelEN).(string),
			Description: util.TernaryOperator(isVI, typeParent.DescVI, typeParent.DescEN).(string),
			Children:    []*AccommodationTypeChildren{},
		}
		for _, typeChild := range typeChilds {
			if typeParent.Class == typeChild.Class {
				typeResp.Children = append(typeResp.Children, &AccommodationTypeChildren{
					ID:          typeChild.ID,
					Label:       util.TernaryOperator(isVI, typeChild.LabelVI, typeChild.LabelEN).(string),
					Description: util.TernaryOperator(isVI, typeChild.DescVI, typeChild.DescEN).(string),
				})
			}
		}

		resp = append(resp, typeResp)
	}

	return resp, nil
}

func (s *GeneralType) ListBed(ctx context.Context, lang string) ([]*BedTypeResponse, error) {

	bedTypes, err := s.getTypeByOrder(ctx, []string{
		model.GeneralTypeClassBed,
	}, lang, 1)
	if err != nil {
		return nil, err
	}

	resp := []*BedTypeResponse{}
	isVI := lang == model.GeneralTypeLangVI
	for _, bedType := range bedTypes {
		resp = append(resp, &BedTypeResponse{
			ID:    bedType.ID,
			Label: util.TernaryOperator(isVI, bedType.LabelVI, bedType.LabelEN).(string),
		})
	}

	return resp, nil
}
