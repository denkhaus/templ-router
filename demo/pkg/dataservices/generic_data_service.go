package dataservices

import (
	"context"
	"fmt"

	"github.com/samber/do/v2"
)

// ProductData represents product information
type ProductData struct {
	ID          string
	Name        string
	Price       float64
	Description string
	Category    string
	InStock     bool
}

// OrderData represents order information
type OrderData struct {
	ID       string
	UserID   string
	Products []ProductData
	Total    float64
	Status   string
}

// GenericDataService kann verschiedene Datentypen zurückgeben
// basierend auf dem "type" Parameter in der Route
type GenericDataService interface {
	GetUserData(ctx context.Context, params map[string]string) (*UserData, error)
	GetProductData(ctx context.Context, params map[string]string) (*ProductData, error)
	GetOrderData(ctx context.Context, params map[string]string) (*OrderData, error)
}

// genericDataServiceImpl implementiert alle Datentypen
type genericDataServiceImpl struct {
	// Hier würdest du normalerweise DB-Connections, APIs, etc. haben
}

// NewGenericDataService erstellt einen neuen generischen DataService
func NewGenericDataService(injector do.Injector) (GenericDataService, error) {
	return &genericDataServiceImpl{}, nil
}

// GetUserData implementiert UserDataService Interface
func (s *genericDataServiceImpl) GetUserData(ctx context.Context, params map[string]string) (*UserData, error) {
	userID := params["id"]
	if userID == "" {
		userID = "demo-user"
	}

	return &UserData{
		ID:       userID,
		Name:     fmt.Sprintf("User %s", userID),
		Email:    fmt.Sprintf("%s@example.com", userID),
		Role:     "User",
		Projects: 5,
		Tasks:    12,
	}, nil
}

// GetProductData für Product-Templates
func (s *genericDataServiceImpl) GetProductData(ctx context.Context, params map[string]string) (*ProductData, error) {
	productID := params["id"]
	if productID == "" {
		productID = "demo-product"
	}

	return &ProductData{
		ID:          productID,
		Name:        fmt.Sprintf("Product %s", productID),
		Price:       99.99,
		Description: "This is a demo product",
		Category:    "Electronics",
		InStock:     true,
	}, nil
}

// GetOrderData für Order-Templates
func (s *genericDataServiceImpl) GetOrderData(ctx context.Context, params map[string]string) (*OrderData, error) {
	orderID := params["id"]
	if orderID == "" {
		orderID = "demo-order"
	}

	return &OrderData{
		ID:     orderID,
		UserID: "user-123",
		Products: []ProductData{
			{ID: "prod-1", Name: "Product 1", Price: 29.99},
			{ID: "prod-2", Name: "Product 2", Price: 49.99},
		},
		Total:  79.98,
		Status: "Shipped",
	}, nil
}