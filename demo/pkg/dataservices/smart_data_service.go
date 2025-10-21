package dataservices

import (
	"context"
	"fmt"
	"reflect"

	"github.com/samber/do/v2"
)

// SmartDataService kann automatisch den richtigen Typ zurückgeben
// basierend auf dem erwarteten Return-Type des Templates
type SmartDataService interface {
	// GetData ist die universelle Methode - der Return-Type wird zur Laufzeit bestimmt
	GetData(ctx context.Context, params map[string]string) (interface{}, error)
}

// smartDataServiceImpl implementiert die intelligente Datenauflösung
type smartDataServiceImpl struct {
	// Hier könntest du verschiedene Datenquellen haben
}

// NewSmartDataService erstellt einen neuen intelligenten DataService
func NewSmartDataService(injector do.Injector) (SmartDataService, error) {
	return &smartDataServiceImpl{}, nil
}

// GetData bestimmt automatisch welcher Datentyp zurückgegeben werden soll
// basierend auf dem "dataType" Parameter oder der Route
func (s *smartDataServiceImpl) GetData(ctx context.Context, params map[string]string) (interface{}, error) {
	// Option 1: Expliziter dataType Parameter
	dataType := params["dataType"]
	
	// Option 2: Ableitung vom Route-Pattern
	if dataType == "" {
		// Du könntest den Typ aus der Route ableiten
		if userID := params["userId"]; userID != "" {
			dataType = "user"
		} else if productID := params["productId"]; productID != "" {
			dataType = "product"
		} else if orderID := params["orderId"]; orderID != "" {
			dataType = "order"
		} else {
			// Fallback basierend auf Route-Kontext
			dataType = "user" // Default
		}
	}

	switch dataType {
	case "user":
		return s.getUserData(ctx, params)
	case "product":
		return s.getProductData(ctx, params)
	case "order":
		return s.getOrderData(ctx, params)
	default:
		return nil, fmt.Errorf("unknown data type: %s", dataType)
	}
}

// Private Methoden für spezifische Datentypen
func (s *smartDataServiceImpl) getUserData(ctx context.Context, params map[string]string) (*UserData, error) {
	userID := params["id"]
	if userID == "" {
		userID = "smart-user"
	}

	return &UserData{
		ID:       userID,
		Name:     fmt.Sprintf("Smart User %s", userID),
		Email:    fmt.Sprintf("%s@smart.com", userID),
		Role:     "Smart User",
		Projects: 8,
		Tasks:    20,
	}, nil
}

func (s *smartDataServiceImpl) getProductData(ctx context.Context, params map[string]string) (*ProductData, error) {
	productID := params["id"]
	if productID == "" {
		productID = "smart-product"
	}

	return &ProductData{
		ID:          productID,
		Name:        fmt.Sprintf("Smart Product %s", productID),
		Price:       149.99,
		Description: "This is a smart product with AI features",
		Category:    "Smart Electronics",
		InStock:     true,
	}, nil
}

func (s *smartDataServiceImpl) getOrderData(ctx context.Context, params map[string]string) (*OrderData, error) {
	orderID := params["id"]
	if orderID == "" {
		orderID = "smart-order"
	}

	return &OrderData{
		ID:     orderID,
		UserID: "smart-user-123",
		Products: []ProductData{
			{ID: "smart-1", Name: "Smart Device 1", Price: 99.99},
			{ID: "smart-2", Name: "Smart Device 2", Price: 149.99},
		},
		Total:  249.98,
		Status: "Processing",
	}, nil
}

// TypedDataService ist eine erweiterte Version die Reflection nutzt
type TypedDataService interface {
	GetDataOfType(ctx context.Context, params map[string]string, targetType reflect.Type) (interface{}, error)
}

// Implementierung mit Reflection für maximale Flexibilität
func (s *smartDataServiceImpl) GetDataOfType(ctx context.Context, params map[string]string, targetType reflect.Type) (interface{}, error) {
	// Bestimme den Typ basierend auf dem erwarteten Return-Type
	switch targetType.String() {
	case "*dataservices.UserData":
		return s.getUserData(ctx, params)
	case "*dataservices.ProductData":
		return s.getProductData(ctx, params)
	case "*dataservices.OrderData":
		return s.getOrderData(ctx, params)
	default:
		return nil, fmt.Errorf("unsupported target type: %s", targetType.String())
	}
}