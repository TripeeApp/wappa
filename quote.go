package wappa

import (
	"context"
	"net/http"
	"time"
)

const quoteEndpoint endpoint = `estimate`

var quoteFields = map[string]string{
	"placeOrigin": "PlaceIdOrigin",
	"placeDestiny": "PlaceIdDestiny",
	"latOrigin": "LatitudeOrigin",
	"lngOrigin": "LongitudeOrigin",
	"latDest": "LatitudeDestiny",
	"lngDest": "LongitudeDestiny",
	"employee": "EmployeeId",
}

type Estimate struct {
	Minimum float64 `json:"minimum"`
	Maximum float64 `json:"maximum"`
	Distance float64 `json:"distance"`
	Journey float64 `json:"distance"`
}

type Icon struct {
	Default string `json:"default"`
	OnFocus string `json:"onFocus"`
	Pin string `json:"pin"`
}

// The available subcategories for this category.
type SubCategory struct {
	ID int `json:"id"`
	TypeID int `json:"typeId"`
	Default bool `json:"default"`
	Description string `json:"description"`
	Discount float64 `json:"discount"`
	Estimate Estimate  `json:"estimate"`
	Observation string `json:"observation"`
	Icon Icon  `json:"icon"`
}

// The categories available for the region requested.
type Category struct {
	ID int `json:"id"`
	Description string `json:"description"`
	SubCategories []SubCategory  `json:"subcategories"`
}

// DriverResult is the API response payload.
type QuoteResult struct {
	Categories []*Category  `json:"categories"`
	EstimatedAt *time.Time `json:"date"`
}

// QuoteService is responsible for handling
// the requests to the quote resource.
type QuoteService service

// Estimate returns a quote for a ride with the given parameters.
func (qs *QuoteService) Estimate(ctx context.Context, f Filter) (*QuoteResult, error) {
	q := &QuoteResult{}

	if err := qs.client.Request(ctx, http.MethodGet, quoteEndpoint.Query(f.Values(quoteFields)), nil, q); err != nil {
		return nil, err
	}

	return q, nil
}
