package dataservices

import (
	"github.com/denkhaus/templ-router/pkg/interfaces"
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
	Locale   string // URL parameter field
	// Query parameter fields
	Page     string
	PageSize string
	Filter   string
	Sort     string
}

// UserDataService provides user data for templates
// Uses standardized GetData method for all DataServices
type UserDataService interface {
	GetData(routerCtx interfaces.RouterContext) (*UserData, error)
	GetUserData(routerCtx interfaces.RouterContext) (*UserData, error)
}

// userDataServiceImpl is the concrete implementation
type userDataServiceImpl struct {
	// In a real app, this would have DB connections, etc.
}

// NewUserDataService creates a new user data service
func NewUserDataService(injector do.Injector) (UserDataService, error) {
	return &userDataServiceImpl{}, nil
}

// GetData retrieves user data based on route parameters and query parameters
func (s *userDataServiceImpl) GetData(routerCtx interfaces.RouterContext) (*UserData, error) {
	// URL Parameters
	locale := routerCtx.GetURLParam("locale")
	if locale == "" {
		locale = "en"
	}

	// Query Parameters - NEW: RouterContext query parameter support!
	page := routerCtx.GetQueryParam("page")
	pageSize := routerCtx.GetQueryParam("pageSize")
	filter := routerCtx.GetQueryParam("filter")
	sort := routerCtx.GetQueryParam("sort")

	// Set defaults for query parameters
	if page == "" {
		page = "1"
	}
	if pageSize == "" {
		pageSize = "10"
	}
	if sort == "" {
		sort = "name"
	}

	// Return demo data with query parameter information
	userData := &UserData{
		ID:       "demo-test-user",
		Name:     "DataService Demo User",
		Email:    "demo@example.com",
		Role:     "Test User",
		Projects: 3,
		Tasks:    15,
		// Query parameter data
		Page:     page,
		PageSize: pageSize,
		Filter:   filter,
		Sort:     sort,
	}
	
	// Add locale to UserData for display
	userData.Locale = locale

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

// GetUserData is the specific method that should be called preferentially over GetData
func (s *userDataServiceImpl) GetUserData(routerCtx interfaces.RouterContext) (*UserData, error) {
	// URL Parameters
	locale := routerCtx.GetURLParam("locale")
	if locale == "" {
		locale = "en"
	}

	// Query Parameters - NEW: RouterContext query parameter support!
	page := routerCtx.GetQueryParam("page")
	pageSize := routerCtx.GetQueryParam("pageSize")
	filter := routerCtx.GetQueryParam("filter")
	sort := routerCtx.GetQueryParam("sort")

	// Set defaults for query parameters
	if page == "" {
		page = "1"
	}
	if pageSize == "" {
		pageSize = "10"
	}
	if sort == "" {
		sort = "name"
	}

	// Return demo data with a marker to show GetUserData method was called
	userData := &UserData{
		ID:       "specific-method-user",
		Name:     "GetUserData Method Called!",
		Email:    "getuserdata@example.com",
		Role:     "Specific Method User",
		Projects: 5,
		Tasks:    25,
		Locale:   locale, // URL parameter
		// Query parameter data
		Page:     page,
		PageSize: pageSize,
		Filter:   filter,
		Sort:     sort,
	}

	// Add some variation based on locale
	switch locale {
	case "de":
		userData.Name = "GetUserData Methode aufgerufen!"
		userData.Email = "getuserdata@beispiel.de"
		userData.Role = "Spezifische Methode Benutzer"
	case "es":
		userData.Name = "Metodo GetUserData llamado!"
		userData.Email = "getuserdata@ejemplo.es"
		userData.Role = "Usuario Metodo Especifico"
	case "fr":
		userData.Name = "Methode GetUserData appelee!"
		userData.Email = "getuserdata@exemple.fr"
		userData.Role = "Utilisateur Methode Specifique"
	}

	return userData, nil
}
