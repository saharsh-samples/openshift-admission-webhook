package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/saharsh-samples/go-mux-sql-starter/test"
)

func TestMustLoadEnvValueAsString_when_value_is_set(t *testing.T) {

	// arrange
	expectedValue := "expected-value"
	os.Setenv("ENV_KEY", expectedValue)

	// act
	value := MustLoadEnvValueAsString("ENV_KEY")

	// assert
	test.AssertEquals("", expectedValue, value, t)
}

func TestMustLoadEnvValueAsString_when_value_is_not_set(t *testing.T) {

	// arrange
	os.Setenv("ENV_KEY", "")

	// assert panic
	defer test.AssertPanic("Expected missing environment variable to cause panic", t)

	// act
	MustLoadEnvValueAsString("ENV_KEY")
}

func TestLoadEnvValueAsString_when_value_is_set(t *testing.T) {

	// arrange
	expectedValue := "expected-value"
	defaultValue := "default-value"
	os.Setenv("ENV_KEY", expectedValue)

	// act
	value := LoadEnvValueAsString("ENV_KEY", defaultValue)

	// assert
	test.AssertEquals("", expectedValue, value, t)
}

func TestLoadEnvValueAsString_when_value_is_not_set(t *testing.T) {

	// arrange
	defaultValue := "default-value"
	os.Setenv("ENV_KEY", "")

	// act
	value := LoadEnvValueAsString("ENV_KEY", defaultValue)

	// assert
	test.AssertEquals("", defaultValue, value, t)
}

func TestLoadEnvValueAsInteger_when_value_is_set(t *testing.T) {

	// arrange
	expectedValue := 5
	defaultValue := 3
	os.Setenv("ENV_KEY", fmt.Sprintf("%d", expectedValue))

	// act
	value := LoadEnvValueAsInteger("ENV_KEY", defaultValue)

	// assert
	test.AssertEquals("", expectedValue, value, t)
}

func TestLoadEnvValueAsInteger_when_value_is_not_set(t *testing.T) {

	// arrange
	defaultValue := 3
	os.Setenv("ENV_KEY", "")

	// act
	value := LoadEnvValueAsInteger("ENV_KEY", defaultValue)

	// assert
	test.AssertEquals("", defaultValue, value, t)
}

func TestLoadEnvValueAsInteger_when_value_is_not_integer(t *testing.T) {

	// arrange
	os.Setenv("ENV_KEY", "hi")

	// assert via defer
	defer test.AssertPanic("Expected panic", t)

	// act
	LoadEnvValueAsInteger("ENV_KEY", 0)
}

func TestLoadEnvValueAsBool_when_value_is_set_to_true(t *testing.T) {

	// arrange
	os.Setenv("ENV_KEY", "true")

	// act
	value := LoadEnvValueAsBool("ENV_KEY", false)

	// assert
	test.AssertTrue("Expected 'true'", value, t)
}

func TestLoadEnvValueAsBool_when_value_is_set_to_True(t *testing.T) {

	// arrange
	os.Setenv("ENV_KEY", "True")

	// act
	value := LoadEnvValueAsBool("ENV_KEY", false)

	// assert
	test.AssertTrue("Expected 'true'", value, t)
}

func TestLoadEnvValueAsBool_when_value_is_set_to_tRuE(t *testing.T) {

	// arrange
	os.Setenv("ENV_KEY", "tRuE")

	// act
	value := LoadEnvValueAsBool("ENV_KEY", false)

	// assert
	test.AssertTrue("Expected 'true'", value, t)
}

func TestLoadEnvValueAsBool_when_value_is_set_to_false(t *testing.T) {

	// arrange
	os.Setenv("ENV_KEY", "false")

	// act
	value := LoadEnvValueAsBool("ENV_KEY", true)

	// assert
	test.AssertFalse("Expected 'false'", value, t)
}

func TestLoadEnvValueAsBool_when_value_is_set_to_False(t *testing.T) {

	// arrange
	os.Setenv("ENV_KEY", "False")

	// act
	value := LoadEnvValueAsBool("ENV_KEY", true)

	// assert
	test.AssertFalse("Expected 'false'", value, t)
}

func TestLoadEnvValueAsBool_when_value_is_set_to_FaLsE(t *testing.T) {

	// arrange
	os.Setenv("ENV_KEY", "FaLsE")

	// act
	value := LoadEnvValueAsBool("ENV_KEY", true)

	// assert
	test.AssertFalse("Expected 'false'", value, t)
}

func TestLoadEnvValueAsBool_when_value_is_not_set(t *testing.T) {

	// arrange
	os.Setenv("ENV_KEY", "")

	// act
	resultWithTrueDefault := LoadEnvValueAsBool("ENV_KEY", true)
	resultWithFalseDefault := LoadEnvValueAsBool("ENV_KEY", false)

	// assert
	test.AssertTrue("Expected 'true'", resultWithTrueDefault, t)
	test.AssertFalse("Expected 'false'", resultWithFalseDefault, t)
}

func TestLoadEnvValueAsInteger_when_value_is_not_bool(t *testing.T) {

	// arrange
	os.Setenv("ENV_KEY", "hi")

	// assert via defer
	defer test.AssertPanic("Expected panic", t)

	// act
	LoadEnvValueAsBool("ENV_KEY", true)
}

func TestInitFromEnvironment(t *testing.T) {

	// arrange
	os.Setenv("MAW_SERVER_STARTUP_TIMEOUT_SECONDS", "10")
	os.Setenv("PORT", "80")
	os.Setenv("MAW_TLS_CERT_FILE_PATH", "tls-cert-path")
	os.Setenv("MAW_TLS_KEY_FILE_PATH", "tls-key-path")

	// act
	appConfig := &AppConfig{}
	appConfig.InitFromEnvironment()

	// assert
	test.AssertEquals("", 10, appConfig.HTTPConfig.StartupTimeoutInSeconds, t)
	test.AssertEquals("", 80, appConfig.HTTPConfig.Port, t)
	test.AssertEquals("", "tls-cert-path", appConfig.HTTPConfig.TLSCertFilePath, t)
	test.AssertEquals("", "tls-key-path", appConfig.HTTPConfig.TLSKeyFilePath, t)

}

func TestInitFromEnvironment_for_defaults(t *testing.T) {

	// arrange
	os.Setenv("MAW_SERVER_STARTUP_TIMEOUT_SECONDS", "")
	os.Setenv("PORT", "")
	os.Setenv("MAW_TLS_CERT_FILE_PATH", "")
	os.Setenv("MAW_TLS_KEY_FILE_PATH", "")

	// act
	appConfig := &AppConfig{}
	appConfig.InitFromEnvironment()

	// assert
	test.AssertEquals("", 5, appConfig.HTTPConfig.StartupTimeoutInSeconds, t)
	test.AssertEquals("", 8443, appConfig.HTTPConfig.Port, t)
	test.AssertEquals("", "/tls/tls.crt", appConfig.HTTPConfig.TLSCertFilePath, t)
	test.AssertEquals("", "/tls/tls.key", appConfig.HTTPConfig.TLSKeyFilePath, t)
}
