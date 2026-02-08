package repomodel

import "time"

type Part struct {
	// Unique identifier of the part.
	Uuid string
	// Name of the part.
	Name string
	// Description of the part.
	Description string
	// Unit price.
	PriceMinor int64
	// Quantity available in stock.
	StockQuantity int64
	// Part category.
	Category Category
	// Part dimensions.
	Dimensions *Dimensions
	// Manufacturer information.
	Manufacturer *Manufacturer
	// Tags for quick search.
	Tags []string
	// Flexible metadata.
	Metadata map[string]*Value
	// Creation timestamp.
	CreatedAt *time.Time
	// Last update timestamp.
	UpdatedAt *time.Time
}

// Category of the Part.
type Category int32

const (
	CategoryUnspecified Category = 0
	CategoryEngine      Category = 1
	CategoryFuel        Category = 2
	CategoryPorthole    Category = 3
	CategoryWing        Category = 4
)

// Dimenstions of the Part.
type Dimensions struct {
	Length float64
	Width  float64
	Height float64
	Weight float64
}

// Manufacturer of the Part.
type Manufacturer struct {
	Name    string
	Country string
	Website string
}

// Value represents a typed metadata value.
type Value struct {
	StringValue *string
	Int64Value  *int64
	DoubleValue *float64
	BoolValue   *bool
}

type PartsFilter struct {
	Uuids                 []string
	Names                 []string
	Categories            []Category
	ManufacturerCountries []string
	Tags                  []string
}
