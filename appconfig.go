package appconfig

// ActiveProfileEnvName represents environment name for active profile
var ActiveProfileEnvName = "app.profiles.active"

// ConfigFilenamePrefix represents configuration filename prefix
var ConfigFilenamePrefix = "application"

// LoadConfig loads configuration
func LoadConfig(configFilename string, to interface{}) error {
	return nil
}
