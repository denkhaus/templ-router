package dataservices

import (
	"context"

	"github.com/samber/do/v2"
)

// SpecificData represents data that is only available through specific method
type UserWithIdData struct {
	ID            string
	CurrentLocale string
}

// SpecificOnlyDataService provides data ONLY through GetSpecificData method
// This service intentionally does NOT implement GetData method
type UserWithIdDataService interface {
	// NO GetData method here - only the specific method
	GetUserWithIdData(ctx context.Context, params map[string]string) (*UserWithIdData, error)
}

// uaerWithIdDataServiceImpl is the concrete implementation
type userWithIdDataServiceImpl struct{}

// NewSpecificOnlyDataService creates a new specific-only data service
func NewUserWithIdDataService(injector do.Injector) (UserWithIdDataService, error) {
	return &userWithIdDataServiceImpl{}, nil
}

// GetSpecificData retrieves data - this is the ONLY method available
// No GetData method is implemented
func (s *userWithIdDataServiceImpl) GetUserWithIdData(ctx context.Context, params map[string]string) (*UserWithIdData, error) {
	locale := params["locale"]
	if locale == "" {
		locale = "undefined"
	}

	userId := params["userId"]
	if userId == "" {
		userId = "undefined"
	}

	// Return demo data showing that the specific method was called
	data := &UserWithIdData{
		ID:            userId,
		CurrentLocale: locale,
	}

	return data, nil
}
