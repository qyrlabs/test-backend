package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/qyrlabs/test-backend/inventory/internal/model"
	inventoryv1 "github.com/qyrlabs/test-backend/shared/pkg/proto/inventory/v1"
)

func ToProtoPart(part *model.Part) *inventoryv1.Part {
	return &inventoryv1.Part{
		Uuid:          part.Uuid,
		Name:          part.Name,
		Description:   part.Description,
		PriceMinor:    part.PriceMinor,
		StockQuantity: part.StockQuantity,
		Category:      ToProtoCategory(part.Category),
		Dimensions:    ToProtoDimensions(part.Dimensions),
		Manufacturer:  ToProtoManufacturer(part.Manufacturer),
		Tags:          part.Tags,
		Metadata:      ToProtoValueMap(part.Metadata),
		CreatedAt:     timestamppb.New(*part.CreatedAt),
		UpdatedAt:     timestamppb.New(*part.UpdatedAt),
	}
}

func ToProtoParts(parts []*model.Part) []*inventoryv1.Part {
	protoParts := make([]*inventoryv1.Part, 0, len(parts))
	for _, part := range parts {
		protoParts = append(protoParts, ToProtoPart(part))
	}
	return protoParts
}

func ToProtoCategory(category model.Category) inventoryv1.Category {
	switch category {
	case model.CategoryUnspecified:
		return inventoryv1.Category_CATEGORY_UNSPECIFIED
	case model.CategoryEngine:
		return inventoryv1.Category_CATEGORY_ENGINE
	case model.CategoryFuel:
		return inventoryv1.Category_CATEGORY_FUEL
	case model.CategoryPorthole:
		return inventoryv1.Category_CATEGORY_PORTHOLE
	case model.CategoryWing:
		return inventoryv1.Category_CATEGORY_WING
	default:
		return inventoryv1.Category_CATEGORY_UNSPECIFIED
	}
}

func ToModelCategory(category inventoryv1.Category) model.Category {
	switch category {
	case inventoryv1.Category_CATEGORY_UNSPECIFIED:
		return model.CategoryUnspecified
	case inventoryv1.Category_CATEGORY_ENGINE:
		return model.CategoryEngine
	case inventoryv1.Category_CATEGORY_FUEL:
		return model.CategoryFuel
	case inventoryv1.Category_CATEGORY_PORTHOLE:
		return model.CategoryPorthole
	case inventoryv1.Category_CATEGORY_WING:
		return model.CategoryWing
	default:
		return model.CategoryUnspecified
	}
}

func ToProtoDimensions(dimensions *model.Dimensions) *inventoryv1.Dimensions {
	if dimensions == nil {
		return nil
	}

	return &inventoryv1.Dimensions{
		Length: dimensions.Length,
		Width:  dimensions.Width,
		Height: dimensions.Height,
		Weight: dimensions.Weight,
	}
}

func ToProtoManufacturer(manufacturer *model.Manufacturer) *inventoryv1.Manufacturer {
	if manufacturer == nil {
		return nil
	}

	return &inventoryv1.Manufacturer{
		Name:    manufacturer.Name,
		Country: manufacturer.Country,
		Website: manufacturer.Website,
	}
}

func ToProtoValueMap(metadata map[string]*model.Value) map[string]*inventoryv1.Value {
	if metadata == nil {
		return nil
	}

	res := make(map[string]*inventoryv1.Value, len(metadata))
	for key, value := range metadata {
		res[key] = ToProtoValue(value)
	}

	return res
}

func ToProtoValue(value *model.Value) *inventoryv1.Value {
	if value == nil {
		return nil
	}

	protoValue := &inventoryv1.Value{}

	switch {
	case value.StringValue != nil:
		protoValue.Kind = &inventoryv1.Value_StringValue{StringValue: *value.StringValue}
	case value.Int64Value != nil:
		protoValue.Kind = &inventoryv1.Value_Int64Value{Int64Value: *value.Int64Value}
	case value.DoubleValue != nil:
		protoValue.Kind = &inventoryv1.Value_DoubleValue{DoubleValue: *value.DoubleValue}
	case value.BoolValue != nil:
		protoValue.Kind = &inventoryv1.Value_BoolValue{BoolValue: *value.BoolValue}
	default:
		return nil
	}

	return protoValue
}

func ToProtoFilter(filter *inventoryv1.PartsFilter) model.PartsFilter {
	if filter == nil {
		return model.PartsFilter{}
	}

	categories := make([]model.Category, 0, len(filter.GetCategories()))
	for _, cat := range filter.GetCategories() {
		categories = append(categories, ToModelCategory(cat))
	}

	return model.PartsFilter{
		Uuids:                 copyPartsFilterField(filter.GetUuids()),
		Names:                 copyPartsFilterField(filter.GetNames()),
		Categories:            categories,
		ManufacturerCountries: copyPartsFilterField(filter.GetManufacturerCountries()),
		Tags:                  copyPartsFilterField(filter.GetTags()),
	}
}

func copyPartsFilterField(v []string) []string {
	if len(v) == 0 {
		return nil
	}
	res := make([]string, len(v))
	copy(res, v)
	return res
}
