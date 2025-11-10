package config

// IsProduction returns true if running in production mode
func (c *configImpl) IsProduction() bool {
	if c.Environment.Kind == "develop" {
		return false
	}

	return true
}

// IsDevelopment returns true if running in development mode
func (c *configImpl) IsDevelopment() bool {
	return !c.IsProduction()
}
