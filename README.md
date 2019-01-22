[![Build Status](https://travis-ci.com/lechuckroh/appconfig.svg?branch=master)](https://travis-ci.com/lechuckroh/appconfig)

# appconfig
SpringBoot like profile based application configuration loader.

See [confita](https://github.com/heetch/confita) for details.

## Settings
#### `ActiveProfileEnvName`
Environment variable name to specify active profile.

Default value is `app.profiles.active`

```go
appconfig.ActiveProfileEnvName = "app.profiles.active"
```

#### `ConfigFilenamePrefix`
Configuration filename prefix.

Default value is `application`.

```go
appconfig.ConfigFilenamePrefix = "application"
```

## Load Precedence
1. Environment variable.
    * Environment variable is read from `config` tag.
    * If field tag is `nested.first-name`, field value will be read from `NESTED.FIRST_NAME` environment variable if set.
    * List items are separated by comma(`,`).
1. Specified configFilename file (passed as `LoadConfig()` argument)
1. `config/application-{profile}.yaml`
1. `application-{profile}.yaml`
1. `config/application.yaml`
1. `application.yaml`
1. `config/application-{profile}.yaml`

## Example
Suppose we have set environment variables as:
* `FOO_CONFIG_FILE`=`myconfig.yml`
* `foo.profiles.active`=`dev`

The following code will load files in the following order:
1. `myconfig.yml`
1. `config/foo-dev.yaml`
1. `foo-dev.yaml`
1. `config/foo.yaml`
1. `foo.yaml`

```go
package appconfig

import (
	"fmt"
	"github.com/lechuckroh/appconfig"
	"os"
)

type Config struct {
	Name   string   `config:"name"`
	Age    int      `config:"age"`
	Tags   []string `config:"tags"`
	Nested struct {
		Name string `config:"nested.name"`
		Age  int    `config:"nested.age"`
	}
}

func main() {
	appconfig.ActiveProfileEnvName = "foo.profiles.active"
	appconfig.ConfigFilenamePrefix = "foo"
	configFilename := os.Getenv("FOO_CONFIG_FILE")

	config := Config{}
	loadedFilenames, err := appconfig.LoadConfig(configFilename, &config)
	if err != nil {
		fmt.Printf("Failed to load config: %s\n", err.Error())
		return
	}
	for idx, filename := range loadedFilenames {
		fmt.Printf("Configuration loaded: [%d] %s\n", idx+1, filename)
	}
	if len(loadedFilenames) == 0 {
		fmt.Println("No config file loaded")
	}

	fmt.Printf("Config: %+v\n", config)
}
```
