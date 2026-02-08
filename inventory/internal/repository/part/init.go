package part

import (
	"fmt"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/samber/lo"

	"github.com/qyrlabs/test-backend/inventory/internal/repository/repomodel"
)

func (r *repository) initParts(count int) error {
	if count < 0 {
		return fmt.Errorf("failed to execute initParts")
	}

	parts := createParts(count)
	for _, part := range parts {
		r.parts[part.Uuid] = part
	}
	return nil
}

func fakeDimensions() *repomodel.Dimensions {
	return &repomodel.Dimensions{
		Length: gofakeit.Float64Range(1.0, 300.0),
		Width:  gofakeit.Float64Range(1.0, 300.0),
		Height: gofakeit.Float64Range(0.5, 150.0),
		Weight: gofakeit.Float64Range(0.1, 500.0),
	}
}

func fakeManufacturer() *repomodel.Manufacturer {
	return &repomodel.Manufacturer{
		Name:    gofakeit.Company(),
		Country: gofakeit.Country(),
		Website: gofakeit.URL(),
	}
}

func fakeTags() []string {
	tags := make([]string, 0, 5)
	for range gofakeit.IntRange(1, 5) {
		tags = append(tags, gofakeit.Word())
	}
	return tags
}

func randomCategory() repomodel.Category {
	// Ignore any, but not UNSPECIFIED
	vals := []repomodel.Category{
		repomodel.CategoryEngine,
		repomodel.CategoryFuel,
		repomodel.CategoryPorthole,
		repomodel.CategoryWing,
	}
	return vals[gofakeit.IntRange(0, len(vals)-1)]
}

func createParts(count int) []repomodel.Part {
	parts := make([]repomodel.Part, 0, count)
	for range count {
		parts = append(parts, repomodel.Part{
			Uuid:          gofakeit.UUID(),
			Name:          gofakeit.Name(),
			Description:   gofakeit.Sentence(10),
			PriceMinor:    int64(gofakeit.IntRange(1, 100000)),
			StockQuantity: int64(gofakeit.IntRange(1, 100)),
			Category:      randomCategory(),
			Dimensions:    fakeDimensions(),
			Manufacturer:  fakeManufacturer(),
			Tags:          fakeTags(),
			CreatedAt:     lo.ToPtr(gofakeit.Date()),
			UpdatedAt:     lo.ToPtr(gofakeit.Date()),
		})
	}
	return parts
}
