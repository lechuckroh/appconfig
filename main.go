package appconfig

import (
	"context"
	"fmt"
	"github.com/heetch/confita/backend/env"
	"os"
	"time"

	"github.com/heetch/confita"
	"github.com/heetch/confita/backend"
	"github.com/heetch/confita/backend/file"
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
// 1. Specified 'configFilename' file
//
// 2. 'config/application-{profile}.yaml'
//
// 3. 'application-{profile}.yaml'
//
// 4. 'config/application.yaml'
//
// 5. 'application.yaml'
//
func LoadConfig(configFilename string, to interface{}) ([]string, error) {
	var lookupFiles []string
	var loadFilenames []string

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
	backends := make([]backend.Backend, 0)

	// lower precedence file comes first
	for i := len(lookupFiles)-1; i >= 0; i-- {
		lookupFile := lookupFiles[i]

		if fileExists(lookupFile) {
			backends = append(backends, file.NewBackend(lookupFile))
			loadFilenames = append(loadFilenames, lookupFile)
		}
	}
	backends = append(backends, env.NewBackend())

	// Load configuration from backends
	loader := confita.NewLoader(backends...)

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Second)
	defer cancel()

	return loadFilenames, loader.Load(ctx, to)
}
