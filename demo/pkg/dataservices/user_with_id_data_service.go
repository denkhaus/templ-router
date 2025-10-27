package dataservices

import (
	"github.com/denkhaus/templ-router/pkg/interfaces"
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
	GetUserWithIdData(routerCtx interfaces.RouterContext) (*UserWithIdData, error)
}

// uaerWithIdDataServiceImpl is the concrete implementation
type userWithIdDataServiceImpl struct{}

// NewSpecificOnlyDataService creates a new specific-only data service
func NewUserWithIdDataService(injector do.Injector) (UserWithIdDataService, error) {
	return &userWithIdDataServiceImpl{}, nil
}

// GetSpecificData retrieves data - this is the ONLY method available
// No GetData method is implemented
func (s *userWithIdDataServiceImpl) GetUserWithIdData(routerCtx interfaces.RouterContext) (*UserWithIdData, error) {
	locale := routerCtx.GetURLParam("locale")
	if locale == "" {
		locale = "undefined"
	}

	userId := routerCtx.GetURLParam("userId")
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
