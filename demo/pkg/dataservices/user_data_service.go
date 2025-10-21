package dataservices

import (
	"context"

	"github.com/samber/do/v2"
)

// UserData represents user profile data
type UserData struct {
	ID       string
	Name     string
	Email    string
	Role     string
	Projects int
	Tasks    int
}

// UserDataService provides user data for templates
// Uses standardized GetData method for all DataServices
type UserDataService interface {
	GetData(ctx context.Context, params map[string]string) (*UserData, error)
}

// userDataServiceImpl is the concrete implementation
type userDataServiceImpl struct {
	// In a real app, this would have DB connections, etc.
}

// NewUserDataService creates a new user data service
func NewUserDataService(injector do.Injector) (UserDataService, error) {
	return &userDataServiceImpl{}, nil
}

// GetData retrieves user data based on route parameters
func (s *userDataServiceImpl) GetData(ctx context.Context, params map[string]string) (*UserData, error) {
	// For the test page /{locale}/test, we don't have an ID parameter
	// Instead, we'll return demo data based on locale or just static demo data
	locale := params["locale"]
	if locale == "" {
		locale = "en"
	}

	// Return demo data for the test page
	userData := &UserData{
		ID:       "demo-test-user",
		Name:     "DataService Demo User",
		Email:    "demo@example.com",
		Role:     "Test User",
		Projects: 3,
		Tasks:    15,
	}

	// Add some variation based on locale
	switch locale {
	case "de":
		userData.Name = "DataService Demo Benutzer"
		userData.Email = "demo@beispiel.de"
		userData.Role = "Test Benutzer"
	case "es":
		userData.Name = "Usuario Demo DataService"
		userData.Email = "demo@ejemplo.es"
		userData.Role = "Usuario de Prueba"
	case "fr":
		userData.Name = "Utilisateur Demo DataService"
		userData.Email = "demo@exemple.fr"
		userData.Role = "Utilisateur Test"
	}

	return userData, nil
}
