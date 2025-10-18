package config

import "fmt"

// GetDatabaseDSN returns the database connection string
func (c *configImpl) GetDatabaseDSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.Name,
		c.Database.SSLMode,
	)
}

// GetServerAddress returns the server address
func (c *configImpl) GetServerAddress() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}

// IsProduction returns true if running in production mode
func (c *configImpl) IsProduction() bool {
	if c.Environment.Kind == "develop" {
		return false
	}

	return c.Server.BaseURL != "http://localhost:8080" &&
		c.Security.CSRFSecret != "change-me-in-production"
}

// IsDevelopment returns true if running in development mode
func (c *configImpl) IsDevelopment() bool {
	return !c.IsProduction()
}
