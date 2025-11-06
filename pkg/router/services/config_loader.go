package services

import (
	"fmt"
	"os"
	"strings"

	"github.com/denkhaus/templ-router/pkg/interfaces"
	"github.com/denkhaus/templ-router/pkg/shared"
	"github.com/samber/do/v2"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

// configLoaderImpl implements clean configuration loading
type configLoaderImpl struct {
	logger *zap.Logger
}

// NewConfigLoader creates a new config loader implementation for DI
func NewConfigLoader(i do.Injector) (interfaces.ConfigLoader, error) {
	logger := do.MustInvoke[*zap.Logger](i)
	return &configLoaderImpl{
		logger: logger,
	}, nil
}

// LoadRouteConfig implements router.ConfigLoader
func (cl *configLoaderImpl) LoadRouteConfig(templateFile string) (*shared.ConfigFile, error) {
	return cl.LoadConfig(templateFile)
}

// LoadConfig implements router.ConfigLoader
func (cl *configLoaderImpl) LoadConfig(templatePath string) (*shared.ConfigFile, error) {
	yamlPath := cl.getYAMLPath(templatePath)

	// Check if YAML file exists
	if _, err := os.Stat(yamlPath); os.IsNotExist(err) {
		cl.logger.Debug("No YAML config found for template",
			zap.String("template", templatePath),
			zap.String("yaml_path", yamlPath))
		return nil, nil // No config is not an error
	}

	cl.logger.Debug("Loading config from YAML",
		zap.String("template", templatePath),
		zap.String("yaml_path", yamlPath))

	// Read YAML file
	data, err := os.ReadFile(yamlPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read YAML file %s: %w", yamlPath, err)
	}

	// Parse YAML
	var rawConfig map[string]interface{}
	if err := yaml.Unmarshal(data, &rawConfig); err != nil {
		return nil, fmt.Errorf("failed to parse YAML file %s: %w", yamlPath, err)
	}

	// Convert to shared.ConfigFile
	config := &shared.ConfigFile{}

	// Parse auth settings if present
	if authData, ok := rawConfig["auth"]; ok {
		authSettings, err := cl.parseAuthSettings(authData)
		if err != nil {
			cl.logger.Warn("Failed to parse auth settings",
				zap.String("yaml_path", yamlPath),
				zap.Error(err))
		} else {
			config.Auth = authSettings
		}
	}

	cl.logger.Debug("Config loaded successfully",
		zap.String("template", templatePath),
		zap.Bool("has_auth", config.Auth != nil))

	return config, nil
}

// LoadAuthSettings implements router.ConfigLoader
func (cl *configLoaderImpl) LoadAuthSettings(templatePath string) (*interfaces.AuthSettings, error) {
	config, err := cl.LoadConfig(templatePath)
	if err != nil {
		return nil, err
	}

	if config == nil || config.Auth == nil {
		return nil, nil // No auth settings is not an error
	}

	// Convert *shared.AuthConfig to *interfaces.AuthSettings (they are the same type)
	return (*interfaces.AuthSettings)(config.Auth), nil
}

// getYAMLPath returns the YAML file path for a template
func (cl *configLoaderImpl) getYAMLPath(templatePath string) string {
	if strings.HasSuffix(templatePath, ".templ") {
		return templatePath + ".yaml"
	}
	return templatePath + ".templ.yaml"
}

// parseAuthSettings parses auth settings from YAML data
func (cl *configLoaderImpl) parseAuthSettings(authData interface{}) (*shared.AuthConfig, error) {
	authMap, ok := authData.(map[interface{}]interface{})
	if !ok {
		return nil, fmt.Errorf("auth settings must be a map")
	}

	settings := &shared.AuthConfig{
		Settings: make(map[string]interface{}),
	}

	// Parse auth type
	if typeData, ok := authMap["type"]; ok {
		if typeStr, ok := typeData.(string); ok {
			authType, err := cl.parseAuthType(typeStr)
			if err != nil {
				return nil, fmt.Errorf("invalid auth type: %w", err)
			}
			settings.Type = authType
		}
	}

	// Parse redirect URL
	if redirectData, ok := authMap["redirect_url"]; ok {
		if redirectStr, ok := redirectData.(string); ok {
			settings.RedirectURL = redirectStr
		}
	}

	// Parse roles
	if rolesData, ok := authMap["roles"]; ok {
		if rolesList, ok := rolesData.([]interface{}); ok {
			var roles []string
			for _, role := range rolesList {
				if roleStr, ok := role.(string); ok {
					roles = append(roles, roleStr)
				}
			}
			settings.Roles = roles
		}
	}

	// Parse any additional settings
	for key, value := range authMap {
		keyStr, ok := key.(string)
		if !ok {
			continue
		}
		if keyStr != "type" && keyStr != "redirect_url" && keyStr != "roles" {
			settings.Settings[keyStr] = value
		}
	}

	return settings, nil
}

// parseAuthType converts string to AuthType string
func (cl *configLoaderImpl) parseAuthType(typeStr string) (string, error) {
	switch strings.ToLower(typeStr) {
	case "public", "none":
		return "Public", nil
	case "user", "authenticated", "login", "userrequired":
		return "UserRequired", nil
	case "admin", "administrator", "adminrequired":
		return "AdminRequired", nil
	default:
		return "Public", fmt.Errorf("unknown auth type: %s", typeStr)
	}
}
