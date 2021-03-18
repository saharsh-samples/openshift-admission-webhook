package config

// HTTPConfig contains HTTP related configuration
type HTTPConfig struct {
	StartupTimeoutInSeconds int
	Port                    int
	TLSCertFilePath         string
	TLSKeyFilePath          string
}

// InitFromEnvironment populates configuration from environment variables
func (config *HTTPConfig) InitFromEnvironment() {
	config.StartupTimeoutInSeconds = LoadEnvValueAsInteger("MAW_SERVER_STARTUP_TIMEOUT_SECONDS", 5)
	config.Port = LoadEnvValueAsInteger("PORT", 8443)
	config.TLSCertFilePath = LoadEnvValueAsString("MAW_TLS_CERT_FILE_PATH", "/tls/tls.crt")
	config.TLSKeyFilePath = LoadEnvValueAsString("MAW_TLS_KEY_FILE_PATH", "/tls/tls.key")
}
