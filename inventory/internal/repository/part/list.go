package part

import (
	"context"
	"slices"

	"github.com/qyrlabs/test-backend/inventory/internal/model"
	"github.com/qyrlabs/test-backend/inventory/internal/repository/converter"
)

// Returns List of Parts by filter.
func (r *repository) List(ctx context.Context, filter model.PartsFilter) ([]*model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	uuids := filter.Uuids
	names := filter.Names
	categories := filter.Categories
	countries := filter.ManufacturerCountries
	tags := filter.Tags

	filteredParts := make([]*model.Part, 0)
	for uuid, part := range r.parts {
		if (len(uuids) == 0 || slices.Contains(uuids, uuid)) &&
			(len(names) == 0 || slices.Contains(names, part.Name)) &&
			(len(categories) == 0 || slices.Contains(categories, converter.ToModelCategory(part.Category))) &&
			(len(countries) == 0 || slices.Contains(countries, part.Manufacturer.Country)) &&
			(len(tags) == 0 || slices.Equal(tags, part.Tags)) {
			partModel := converter.ToModelPart(part)
			filteredParts = append(filteredParts, partModel)
		}
	}

	return filteredParts, nil
}
