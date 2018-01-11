package envcfg

import (
	"os"
	"testing"
)

func TestGet(t *testing.T) {
	t.Log("Returns empty string for unknown configs")

	c := &config{}

	if c.Get("unknown") != "" {
		t.Errorf("Unexpected value for unknown config")
	}
}

func TestBundledConfigs(t *testing.T) {
	t.Log("Returns default values for bundled configs")

	config := LoadBundled()

	if config.Get("version") != "detached" {
		t.Errorf("==> Bundled config didn't load default: version")
	}
	if config.Get("logLevel") != "debug" {
		t.Errorf("==> Bundled config didn't load default: logLevel")
	}
	if config.Get("environment") != "dev" {
		t.Errorf("==> Bundled config didn't load default: environment")
	}

	t.Log("Returns environment set values for bundled configs")

	os.Setenv("VERSION", "test")
	config = LoadBundled()

	if config.Get("version") != "test" {
		t.Errorf("==> Bundled config didn't load environment value: version")
	}
	os.Unsetenv("VERSION") // Teardown
}

func TestCustomConfigs(t *testing.T) {
	t.Log("Returns default values for custom configs")

	config := Load(Map{"MY_ENV_VAR": "MY_VALUE"})

	if config.Get("myEnvVar") != "MY_VALUE" {
		t.Errorf("==> Custom config didn't load default: myEnvVar")
	}

	t.Log("Returns environment set values for custom configs")

	os.Setenv("MY_ENV_VAR", "ligeiro")
	config = Load(Map{"MY_ENV_VAR": "DEFAULT_VALUE"})

	if config.Get("myEnvVar") != "ligeiro" {
		t.Errorf("==> Custom config didn't load environment value: '%v'", config.Get("myEnvVar"))
	}
	os.Unsetenv("MY_ENV_VAR") // Teardown

	t.Log("Returns default values for custom configs even when they have same name of bundled ones")

	config = Load(Map{"LOG_LEVEL": "error"})

	if config.Get("logLevel") != "error" {
		t.Errorf("==> Bundled config default took precedence over custom ones: '%v'", config.Get("logLevel"))
	}
}

func TestGetInt(t *testing.T) {
	t.Log("Returns casted int from config")

	config := Load(Map{"MY_ENV_VAR": "30000"})

	if config.GetInt("myEnvVar") != 30000 {
		t.Errorf("==> GetInt config didn't casted value: %v", config.Get("myEnvVar"))
	}
	os.Unsetenv("MY_ENV_VAR") // Teardown
}
