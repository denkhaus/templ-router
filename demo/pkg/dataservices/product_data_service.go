package dataservices

import (
	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/samber/do/v2"
)

// ProductData represents product information
type ProductData struct {
	ID          string
	Name        string
	Description string
	Price       float64
	Category    string
	InStock     bool
}

// ProductDataService provides product data - ONLY has GetData method (no specific method)
type ProductDataService interface {
	GetData(routerCtx interfaces.RouterContext) (*ProductData, error)
}

// productDataServiceImpl is the concrete implementation
type productDataServiceImpl struct{}

// NewProductDataService creates a new product data service
func NewProductDataService(injector do.Injector) (ProductDataService, error) {
	return &productDataServiceImpl{}, nil
}

// GetData retrieves product data based on route parameters
// This service ONLY implements GetData, no specific method like GetProductData
func (s *productDataServiceImpl) GetData(routerCtx interfaces.RouterContext) (*ProductData, error) {
	productID := routerCtx.GetURLParam("id")
	if productID == "" {
		productID = "demo-product"
	}

	// Return demo product data
	product := &ProductData{
		ID:          productID,
		Name:        "Demo Product " + productID,
		Description: "This is a demonstration product loaded via GetData method only",
		Price:       99.99,
		Category:    "Electronics",
		InStock:     true,
	}

	return product, nil
}