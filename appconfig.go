package appconfig

import (
	"context"
	"fmt"
	"github.com/heetch/confita"
	"github.com/heetch/confita/backend"
	"github.com/heetch/confita/backend/env"
	"github.com/heetch/confita/backend/file"
	"os"
	"time"
)

// ActiveProfileEnvName represents environment name for active profile
var ActiveProfileEnvName = "app.profiles.active"

// ConfigFilenamePrefix represents configuration filename prefix
var ConfigFilenamePrefix = "application"

// fileExists check if file with given filename exists.
func fileExists(filename string) bool {
	if _, err := os.Stat(filename); err == nil {
		return true
	} else {
		return false
	}
}

// LoadConfig loads configuration from files or environment variables.
//
// Configuration precedence:
//
// 1. Environment variable
//
// 2. Specified 'configFilename' file
//
// 3. 'config/application-{profile}.yaml'
//
// 4. 'config/application.yaml'
//
// 5. 'application-{profile}.yaml'
//
// 6. 'application.yaml'
//
// 7. 'config/application-{profile}.yaml'
func LoadConfig(configFilename string, to interface{}) error {
	var lookupFiles []string

	// Specific file
	if configFilename != "" {
		lookupFiles = append(lookupFiles, configFilename)
	}

	// Profile specific files
	profile := os.Getenv(ActiveProfileEnvName)
	if profile != "" {
		lookupFiles = append(lookupFiles,
			fmt.Sprintf("config/%v-%v.yml", ConfigFilenamePrefix, profile),
			fmt.Sprintf("%v-%v.yml", ConfigFilenamePrefix, profile),
		)
	}

	// Default files
	lookupFiles = append(lookupFiles,
		fmt.Sprintf("config/%v.yml", ConfigFilenamePrefix),
		fmt.Sprintf("%v.yml", ConfigFilenamePrefix),
	)

	// Backends
	var backends []backend.Backend
	{
		env.NewBackend()
	}

	for _, lookupFile := range lookupFiles {
		if fileExists(lookupFile) {
			backends = append(backends, file.NewBackend(lookupFile))
		}
	}

	// Load configuration from backends
	loader := confita.NewLoader(backends...)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return loader.Load(ctx, to)
}
