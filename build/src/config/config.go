package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// AppConfig describes application's configuration
type AppConfig struct {
	HTTPConfig HTTPConfig
}

// InitFromEnvironment loads app config from environment variables
func (config *AppConfig) InitFromEnvironment() {
	config.HTTPConfig.InitFromEnvironment()
}

// MustLoadEnvValueAsString for specified environment variable. Panics if not set
func MustLoadEnvValueAsString(variableName string) string {
	value := os.Getenv(variableName)
	if value == "" {
		panic(fmt.Sprintf("Required environment variable '%v' not set", variableName))
	}
	return value
}

// LoadEnvValueAsString for specified environment variable. If not set, use default
func LoadEnvValueAsString(variableName string, defaultValue string) string {
	value := os.Getenv(variableName)
	if value == "" {
		return defaultValue
	}
	return value
}

// LoadEnvValueAsInteger for specified environment variable. If not set, use default
func LoadEnvValueAsInteger(variableName string, defaultValue int) int {

	stringValue := os.Getenv(variableName)
	if stringValue == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(stringValue)
	if err != nil {
		fmt.Printf("Error parsing environment value '%v=%v' as an integer\n", variableName, stringValue)
		panic(err)
	}

	return value

}

// LoadEnvValueAsBool for specified environment variable. If not set, use default
func LoadEnvValueAsBool(variableName string, defaultValue bool) bool {

	stringValue := os.Getenv(variableName)
	if stringValue == "" {
		return defaultValue
	}

	caseInsensitiveValue := strings.ToLower(stringValue)
	if caseInsensitiveValue == "true" {
		return true
	}
	if caseInsensitiveValue == "false" {
		return false
	}

	err := fmt.Sprintf("Error parsing environment value '%v=%v' as a bool\n", variableName, stringValue)
	panic(err)

}
