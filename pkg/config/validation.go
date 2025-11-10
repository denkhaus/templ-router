package config

// Validate validates the configuration using centralized validators
func (c *configImpl) Validate() error {
	// Initialize reusable validators

	positiveNumberValidator := NewPositiveNumberValidator("", 1)

	// Validate security configuration
	positiveNumberValidator.fieldName = "security.rate_limit_requests"
	if err := positiveNumberValidator.ValidatePositiveNumber(c.Security.RateLimitRequests); err != nil {
		return err
	}

	return nil
}
