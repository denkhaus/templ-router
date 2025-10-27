package dataservices

import (
	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/samber/do/v2"
)

// SpecificData represents data that is only available through specific method
type SpecificData struct {
	ID          string
	Title       string
	Description string
	Method      string
	Note        string
}

// SpecificOnlyDataService provides data ONLY through GetSpecificData method
// This service intentionally does NOT implement GetData method
type SpecificOnlyDataService interface {
	// NO GetData method here - only the specific method
	GetSpecificData(routerCtx interfaces.RouterContext) (*SpecificData, error)
}

// specificOnlyDataServiceImpl is the concrete implementation
type specificOnlyDataServiceImpl struct{}

// NewSpecificOnlyDataService creates a new specific-only data service
func NewSpecificOnlyDataService(injector do.Injector) (SpecificOnlyDataService, error) {
	return &specificOnlyDataServiceImpl{}, nil
}

// GetSpecificData retrieves data - this is the ONLY method available
// No GetData method is implemented
func (s *specificOnlyDataServiceImpl) GetSpecificData(routerCtx interfaces.RouterContext) (*SpecificData, error) {
	locale := routerCtx.GetURLParam("locale")
	if locale == "" {
		locale = "en"
	}

	// Return demo data showing that the specific method was called
	data := &SpecificData{
		ID:          "specific-only-123",
		Title:       "Specific Method Only Test",
		Description: "This data was loaded using ONLY GetSpecificData method",
		Method:      "GetSpecificData (no GetData fallback available)",
		Note:        "This proves that the router correctly calls the specific method when GetData is not available",
	}

	// Add some variation based on locale
	switch locale {
	case "de":
		data.Title = "Nur Spezifische Methode Test"
		data.Description = "Diese Daten wurden NUR über GetSpecificData Methode geladen"
		data.Method = "GetSpecificData (kein GetData Fallback verfügbar)"
		data.Note = "Dies beweist, dass der Router die spezifische Methode korrekt aufruft, wenn GetData nicht verfügbar ist"
	case "es":
		data.Title = "Prueba Solo Método Específico"
		data.Description = "Estos datos se cargaron usando SOLO el método GetSpecificData"
		data.Method = "GetSpecificData (sin fallback GetData disponible)"
		data.Note = "Esto prueba que el router llama correctamente al método específico cuando GetData no está disponible"
	case "fr":
		data.Title = "Test Méthode Spécifique Seulement"
		data.Description = "Ces données ont été chargées en utilisant SEULEMENT la méthode GetSpecificData"
		data.Method = "GetSpecificData (pas de fallback GetData disponible)"
		data.Note = "Ceci prouve que le routeur appelle correctement la méthode spécifique quand GetData n'est pas disponible"
	}

	return data, nil
}