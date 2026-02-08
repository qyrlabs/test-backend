package converter

import (
	"github.com/qyrlabs/test-backend/inventory/internal/model"
	"github.com/qyrlabs/test-backend/inventory/internal/repository/repomodel"
)

func ToModelPart(part repomodel.Part) *model.Part {
	return &model.Part{
		Uuid:          part.Uuid,
		Name:          part.Name,
		Description:   part.Description,
		PriceMinor:    part.PriceMinor,
		StockQuantity: part.StockQuantity,
		Category:      ToModelCategory(part.Category),
		Dimensions:    ToModelDimensions(part.Dimensions),
		Manufacturer:  ToModelManufacturer(part.Manufacturer),
		Tags:          part.Tags,
		Metadata:      ToModelValueMap(part.Metadata),
		CreatedAt:     part.CreatedAt,
		UpdatedAt:     part.UpdatedAt,
	}
}

func ToModelCategory(category repomodel.Category) model.Category {
	switch category {
	case repomodel.CategoryUnspecified:
		return model.CategoryUnspecified
	case repomodel.CategoryEngine:
		return model.CategoryEngine
	case repomodel.CategoryFuel:
		return model.CategoryFuel
	case repomodel.CategoryPorthole:
		return model.CategoryPorthole
	case repomodel.CategoryWing:
		return model.CategoryWing
	default:
		return model.CategoryUnspecified
	}
}

func ToModelDimensions(dimensions *repomodel.Dimensions) *model.Dimensions {
	if dimensions == nil {
		return nil
	}
	return &model.Dimensions{
		Length: dimensions.Length,
		Width:  dimensions.Width,
		Height: dimensions.Height,
		Weight: dimensions.Weight,
	}
}

func ToModelManufacturer(manufacturer *repomodel.Manufacturer) *model.Manufacturer {
	if manufacturer == nil {
		return nil
	}
	return &model.Manufacturer{
		Name:    manufacturer.Name,
		Country: manufacturer.Country,
		Website: manufacturer.Website,
	}
}

func ToModelValueMap(metadata map[string]*repomodel.Value) map[string]*model.Value {
	if metadata == nil {
		return nil
	}
	result := make(map[string]*model.Value, len(metadata))
	for key, value := range metadata {
		result[key] = ToModelValue(value)
	}
	return result
}

func ToModelValue(value *repomodel.Value) *model.Value {
	if value == nil {
		return nil
	}
	return &model.Value{
		StringValue: value.StringValue,
		Int64Value:  value.Int64Value,
		DoubleValue: value.DoubleValue,
		BoolValue:   value.BoolValue,
	}
}
