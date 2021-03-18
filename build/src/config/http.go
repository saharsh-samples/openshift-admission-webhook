package config

// HTTPConfig contains HTTP related configuration
type HTTPConfig struct {
	StartupTimeoutInSeconds int
	Port                    int
}

// InitFromEnvironment populates configuration from environment variables
func (config *HTTPConfig) InitFromEnvironment() {
	config.StartupTimeoutInSeconds = LoadEnvValueAsInteger("MAW_SERVER_STARTUP_TIMEOUT_SECONDS", 5)
	config.Port = LoadEnvValueAsInteger("PORT", 8080)
}
