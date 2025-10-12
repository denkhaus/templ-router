package services

// ValidationResult contains validation results (consolidated from validation_types.go)
type ValidationResult struct {
	Errors   []ValidationError   `json:"errors"`
	Warnings []ValidationWarning `json:"warnings"`
}

// ValidationError represents a validation error that prevents the application from running
type ValidationError struct {
	Type        string   `json:"type"`
	Message     string   `json:"message"`
	FilePath    string   `json:"file_path"`
	RoutePath   string   `json:"route_path"`
	Suggestions []string `json:"suggestions,omitempty"`
}

// ValidationWarning represents a validation warning that should be addressed but doesn't prevent operation
type ValidationWarning struct {
	Type        string   `json:"type"`
	Message     string   `json:"message"`
	FilePath    string   `json:"file_path"`
	RoutePath   string   `json:"route_path"`
	Suggestions []string `json:"suggestions,omitempty"`
}

// HasErrors returns true if there are validation errors
func (vr *ValidationResult) HasErrors() bool {
	return len(vr.Errors) > 0
}

// HasWarnings returns true if there are validation warnings
func (vr *ValidationResult) HasWarnings() bool {
	return len(vr.Warnings) > 0
}

// AddError adds a validation error to the result
func (vr *ValidationResult) AddError(errorType, message, filePath, routePath string, suggestions ...string) {
	vr.Errors = append(vr.Errors, ValidationError{
		Type:        errorType,
		Message:     message,
		FilePath:    filePath,
		RoutePath:   routePath,
		Suggestions: suggestions,
	})
}

// AddWarning adds a validation warning to the result
func (vr *ValidationResult) AddWarning(errorType, message, filePath, routePath string, suggestions ...string) {
	vr.Warnings = append(vr.Warnings, ValidationWarning{
		Type:        errorType,
		Message:     message,
		FilePath:    filePath,
		RoutePath:   routePath,
		Suggestions: suggestions,
	})
}

// GetErrorCount returns the number of validation errors
func (vr *ValidationResult) GetErrorCount() int {
	return len(vr.Errors)
}

// GetWarningCount returns the number of validation warnings
func (vr *ValidationResult) GetWarningCount() int {
	return len(vr.Warnings)
}

// GetErrorsByType returns all errors of a specific type
func (vr *ValidationResult) GetErrorsByType(errorType string) []ValidationError {
	var errors []ValidationError
	for _, err := range vr.Errors {
		if err.Type == errorType {
			errors = append(errors, err)
		}
	}
	return errors
}

// GetWarningsByType returns all warnings of a specific type
func (vr *ValidationResult) GetWarningsByType(warningType string) []ValidationWarning {
	var warnings []ValidationWarning
	for _, warn := range vr.Warnings {
		if warn.Type == warningType {
			warnings = append(warnings, warn)
		}
	}
	return warnings
}

// Merge combines two validation results
func (vr *ValidationResult) Merge(other *ValidationResult) {
	if other == nil {
		return
	}
	vr.Errors = append(vr.Errors, other.Errors...)
	vr.Warnings = append(vr.Warnings, other.Warnings...)
}