// Load configs from environment handling defaults and expose them.
package envcfg

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// Type alias for environment variables configs, it maps `ENV_VAR` with `default_value`
type Map map[string]string

// Already provided environment variable, custom config overrides them.
var bundledConfigs = Map{
	"ENVIRONMENT": "dev",
	"LOG_LEVEL":   "debug",
	"VERSION":     "detached",
}

// Type that holds parsed config and expose them via Get()
type config struct {
	parsedConfigs Map
}

func (c *config) Get(key string) string {
	if configValue, exists := c.parsedConfigs[key]; exists {
		return configValue
	} else {
		fmt.Println(fmt.Sprintf("Unexpected config, returning empty string: %v", key))
		return ""
	}
}

func (c *config) GetInt64(key string) int64 {
	value, err := strconv.ParseInt(c.Get(key), 10, 64)

	if err != nil {
		log.Fatalf("Key value cannot be casted to int: %v", key)
	}

	return value
}

func (c *config) GetInt(key string) int {
	return int(c.GetInt64(key))
}

func (c *config) GetFloat64(key string) float64 {
	value, err := strconv.ParseFloat(c.Get(key), 64)

	if err != nil {
		log.Fatalf("Key value cannot be casted to float: %v", key)
	}

	return value
}

func (c *config) GetFloat(key string) float32 {
	return float32(c.GetFloat64(key))
}

// Bundled configs (ENVIRONMENT, LOG_LEVEL and VERSION) only, useful for brand new applications that has no extra confs.
func LoadBundled() *config {
	return Load(Map{})
}

// Custom environment config map with bundled configs (first has higher priority). look at example for more details.
// Best way to use is to create a config package with init() func that expose Load() result.
func Load(environmentVariablesWithDefaults Map) *config {
	parsedConfigs := Map{}

	// Load bundledConfigs into custom ones only if custom not define them already
	for environmentVariable, defaultValue := range bundledConfigs {
		if _, exists := environmentVariablesWithDefaults[environmentVariable]; !exists {
			environmentVariablesWithDefaults[environmentVariable] = defaultValue
		}
	}

	// Parse merge custom and bundledConfigs fetching environment variables
	for environmentVariable, defaultValue := range environmentVariablesWithDefaults {
		parsedConfigs[toCamelCase(environmentVariable)] = getEnvVarWithDefault(environmentVariable, defaultValue)
	}

	return &config{parsedConfigs}
}

func getEnvVarWithDefault(env string, defaultValue string) string {
	variable := defaultValue

	if setValue := os.Getenv(env); setValue != "" {
		variable = setValue
	}

	return variable
}

func toCamelCase(s string) string {
	capitalizeNextChar := false
	s = strings.Trim(strings.ToLower(s), " ")
	camelized := ""

	for _, char := range s {
		switch {
		case char >= 'A' && char <= 'Z':
			camelized += string(char)
		case char >= 'a' && char <= 'z':
			if capitalizeNextChar {
				camelized += strings.ToUpper(string(char))
				capitalizeNextChar = false
			} else {
				camelized += string(char)
			}
		case char == '_' || char == ' ' || char == '-':
			capitalizeNextChar = true
		}
	}

	return camelized
}
