package config

import (
	"reflect"
	"testing"

	"os"

	"github.com/stretchr/testify/assert"
)

// TestGetConfig_DefaultPort Verify if GetConfig get the default port
func TestGetConfig_DefaultPort(t *testing.T) {
	//Arrange
	expectedPort := 8001
	os.Setenv("APP_PORT", "")
	appConfig = nil

	//Act
	actualConfig := GetConfig()

	//Assert
	assert.NotNil(t, actualConfig)
	assert.Equal(t, expectedPort, actualConfig.Port)
}

// TestGetConfig_EnvironmentPort verify if GetConfig works with Env variables
func TestGetConfig_EnvironmentPort(t *testing.T) {
	//Arrange
	expectedPort := 8002
	os.Setenv("APP_PORT", "8002")
	appConfig = nil

	//Act
	actualConfig := GetConfig()

	//Assert
	assert.NotNil(t, actualConfig)
	assert.Equal(t, expectedPort, actualConfig.Port)
}

// TestGetConfig_EnvironmentPort verify if GetConfig works with the config file
func TestGetConfig_ConfigFilePort(t *testing.T) {
	//Arrange
	expectedPort := 8001
	os.Setenv("ENVIRONEMENT", "DEV")
	os.Setenv("APP_PORT", "")
	appConfig = nil

	//Act
	actualConfig := GetConfig()

	//Assert
	assert.NotNil(t, actualConfig)
	assert.Equal(t, expectedPort, actualConfig.Port)
}

// TestVerifyAllConfigValues ...
func TestVerifyAllConfigValues(t *testing.T) {
	// run configuration init
	appConfig = nil
	os.Setenv("APP_PORT", "")
	os.Setenv("ENVIRONEMENT", "")

	config := GetConfig()

	// check that all fields were initialised using reflection
	v := reflect.ValueOf(*config)
	for i := 0; i < v.NumField(); i++ {
		// get current field, it's name; value ant type
		f := v.Field(i)
		name := v.Type().Field(i).Name

		// verify that value of this field is not default value
		switch f.Kind() {
		case reflect.Slice: // for slice we compare that it is not empty
			if f.IsNil() {
				t.Errorf("Configuration incorrect. Field %q was not initialized correctly.", name)
			}
		case reflect.String,
			reflect.Int:
			// get default value for these types
			z := reflect.Zero(f.Type())
			if z.Interface() == f.Interface() {
				t.Errorf("Configuration incorrect. Field %q was not initialized correctly. Value \"%v\"", name, f)
			}
		default:
			t.Errorf("Configuration incorrect. Field %q has unknown type %q.", name, f.Type().String())
		}
	}
}
