package dataservices

import (
	"context"

	"github.com/samber/do/v2"
)

// OrderData represents order information
type OrderData struct {
	ID       string
	UserID   string
	Items    []string
	Total    float64
	Status   string
	Date     string
}

// OrderDataService provides order data - has BOTH GetData and GetOrderData methods
type OrderDataService interface {
	GetData(ctx context.Context, params map[string]string) (*OrderData, error)
	GetOrderData(ctx context.Context, params map[string]string) (*OrderData, error)
}

// orderDataServiceImpl is the concrete implementation
type orderDataServiceImpl struct{}

// NewOrderDataService creates a new order data service
func NewOrderDataService(injector do.Injector) (OrderDataService, error) {
	return &orderDataServiceImpl{}, nil
}

// GetData retrieves order data - fallback method
func (s *orderDataServiceImpl) GetData(ctx context.Context, params map[string]string) (*OrderData, error) {
	orderID := params["id"]
	if orderID == "" {
		orderID = "fallback-order"
	}

	return &OrderData{
		ID:     orderID,
		UserID: "user-123",
		Items:  []string{"Item A", "Item B"},
		Total:  149.99,
		Status: "GetData method called (fallback)",
		Date:   "2024-01-01",
	}, nil
}

// GetOrderData retrieves order data - specific method that should be called preferentially
func (s *orderDataServiceImpl) GetOrderData(ctx context.Context, params map[string]string) (*OrderData, error) {
	orderID := params["id"]
	if orderID == "" {
		orderID = "specific-order"
	}

	return &OrderData{
		ID:     orderID,
		UserID: "user-456",
		Items:  []string{"Premium Item X", "Premium Item Y", "Premium Item Z"},
		Total:  299.99,
		Status: "GetOrderData method called (specific method)",
		Date:   "2024-10-26",
	}, nil
}