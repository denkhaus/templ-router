package dataservices

import (
	"context"

	"github.com/samber/do/v2"
)

// BrokenData represents data that will cause errors
type BrokenData struct {
	ID      string
	Message string
}

// BrokenDataService is intentionally broken - has NO GetData or GetBrokenData methods
type BrokenDataService interface {
	// This service intentionally has wrong method names to test error handling
	FetchData(ctx context.Context, params map[string]string) (*BrokenData, error)
	RetrieveData(ctx context.Context, params map[string]string) (*BrokenData, error)
}

// brokenDataServiceImpl is the concrete implementation
type brokenDataServiceImpl struct{}

// NewBrokenDataService creates a new broken data service
func NewBrokenDataService(injector do.Injector) (BrokenDataService, error) {
	return &brokenDataServiceImpl{}, nil
}

// FetchData - wrong method name, should be GetData or GetBrokenData
func (s *brokenDataServiceImpl) FetchData(ctx context.Context, params map[string]string) (*BrokenData, error) {
	return &BrokenData{
		ID:      "broken-1",
		Message: "This method should not be called",
	}, nil
}

// RetrieveData - another wrong method name
func (s *brokenDataServiceImpl) RetrieveData(ctx context.Context, params map[string]string) (*BrokenData, error) {
	return &BrokenData{
		ID:      "broken-2",
		Message: "This method should also not be called",
	}, nil
}