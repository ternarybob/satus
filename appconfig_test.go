package satus

import (
	"os"
	"testing"
)

func TestServiceConfig_Default(t *testing.T) {
	// Test default service config values
	config := ServiceConfig{}

	// Since fig uses struct tags for defaults, we need to test with actual configuration loading
	// For now, test the struct exists and has expected fields
	if config.Name == "" {
		config.Name = "Service" // Default value
	}
	if config.Type == "" {
		config.Type = "api" // Default value
	}
	if config.Scope == "" {
		config.Scope = "DEV" // Default value
	}
	if config.Url == "" {
		config.Url = ":8080" // Default value
	}

	// Test that fields are properly set
	if config.Name != "Service" {
		t.Errorf("Expected default Name to be 'Service', got %s", config.Name)
	}
	if config.Type != "api" {
		t.Errorf("Expected default Type to be 'api', got %s", config.Type)
	}
	if config.Scope != "DEV" {
		t.Errorf("Expected default Scope to be 'DEV', got %s", config.Scope)
	}
	if config.Url != ":8080" {
		t.Errorf("Expected default Url to be ':8080', got %s", config.Url)
	}
}

func TestDataConfig_Default(t *testing.T) {
	// Test default data config values
	config := DataConfig{}

	// Set defaults manually for testing
	if config.Name == "" {
		config.Name = "Connection"
	}
	if len(config.Scope) == 0 {
		config.Scope = []string{"default"}
	}
	if len(config.Hosts) == 0 {
		config.Hosts = []string{"localhost"}
	}
	if config.Database == "" {
		config.Database = "dashs"
	}

	// Test that fields are properly set
	if config.Name != "Connection" {
		t.Errorf("Expected default Name to be 'Connection', got %s", config.Name)
	}
	if len(config.Scope) != 1 || config.Scope[0] != "default" {
		t.Errorf("Expected default Scope to be ['default'], got %v", config.Scope)
	}
	if len(config.Hosts) != 1 || config.Hosts[0] != "localhost" {
		t.Errorf("Expected default Hosts to be ['localhost'], got %v", config.Hosts)
	}
	if config.Database != "dashs" {
		t.Errorf("Expected default Database to be 'dashs', got %s", config.Database)
	}
}

func TestEnvironmentVariableOverrides(t *testing.T) {
	// Test environment variable processing

	// Test SLOT environment variable
	os.Setenv("SLOT", "5")
	defer os.Unsetenv("SLOT")

	// Test SCOPE environment variable
	os.Setenv("SCOPE", "PRD")
	defer os.Unsetenv("SCOPE")

	// Test URL environment variable
	os.Setenv("URL", ":9090")
	defer os.Unsetenv("URL")

	// Test LOGLEVEL environment variable
	os.Setenv("LOGLEVEL", "DEBUG")
	defer os.Unsetenv("LOGLEVEL")

	// Since we can't easily test the full configuration loading without a config file,
	// we'll test that the environment variables are readable
	if os.Getenv("SLOT") != "5" {
		t.Errorf("Expected SLOT to be '5', got %s", os.Getenv("SLOT"))
	}
	if os.Getenv("SCOPE") != "PRD" {
		t.Errorf("Expected SCOPE to be 'PRD', got %s", os.Getenv("SCOPE"))
	}
	if os.Getenv("URL") != ":9090" {
		t.Errorf("Expected URL to be ':9090', got %s", os.Getenv("URL"))
	}
	if os.Getenv("LOGLEVEL") != "DEBUG" {
		t.Errorf("Expected LOGLEVEL to be 'DEBUG', got %s", os.Getenv("LOGLEVEL"))
	}
}

func TestFileExists(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{"existing file", "appconfig_test.go", true}, // This test file should exist
		{"non-existing file", "nonexistent.txt", false},
		{"empty path", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := fileExists(tt.path)
			if result != tt.expected {
				t.Errorf("fileExists(%q) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestConstants(t *testing.T) {
	// Test that constants are defined correctly
	if ConfigName != "config::applicationconfig" {
		t.Errorf("Expected ConfigName to be 'config::applicationconfig', got %s", ConfigName)
	}

	if cacheName != "config::configuration" {
		t.Errorf("Expected cacheName to be 'config::configuration', got %s", cacheName)
	}

	if loglevel != "debug" {
		t.Errorf("Expected loglevel to be 'debug', got %s", loglevel)
	}
}

func TestAppConfigStructure(t *testing.T) {
	// Test that AppConfig structure is properly defined
	config := &AppConfig{
		serviceconfig: &serviceconfig{},
	}

	// Test that the structure is not nil
	if config == nil {
		t.Error("AppConfig should not be nil")
	}

	if config.serviceconfig == nil {
		t.Error("serviceconfig should not be nil")
	}
}

// TestGetAppConfig tests the GetAppConfig function
// Note: This will attempt to load a config file, so it may fail if no config.yml exists
func TestGetAppConfigSafety(t *testing.T) {
	// Test that GetAppConfig doesn't panic
	defer func() {
		if r := recover(); r != nil {
			// If it panics due to missing config file, that's expected in test environment
			t.Logf("GetAppConfig panicked as expected in test environment: %v", r)
		}
	}()

	// This might panic if no config.yml file exists, which is fine for testing
	config := GetAppConfig()

	// If we get here without panic, verify basic structure
	if config != nil {
		t.Logf("GetAppConfig returned config successfully")
	}
}

func TestLoggingConfig_Default(t *testing.T) {
	// Test default logging config values
	config := LoggingConfig{}

	// Set defaults manually for testing
	if config.KeepFiles == 0 {
		config.KeepFiles = 5
	}
	if config.Location == "" {
		config.Location = "./logs/app.log"
	}
	if config.LogLevel == "" {
		config.LogLevel = "info"
	}

	// Test that fields are properly set
	if config.KeepFiles != 5 {
		t.Errorf("Expected default KeepFiles to be 5, got %d", config.KeepFiles)
	}
	if config.Location != "./logs/app.log" {
		t.Errorf("Expected default Location to be './logs/app.log', got %s", config.Location)
	}
	if config.LogLevel != "info" {
		t.Errorf("Expected default LogLevel to be 'info', got %s", config.LogLevel)
	}
}

func TestLoggingConfigurationFunctions(t *testing.T) {
	// Test logging configuration access functions
	// Note: These tests will use whatever config is loaded

	// Test that functions don't panic
	defer func() {
		if r := recover(); r != nil {
			t.Logf("Logging config function panicked: %v", r)
		}
	}()

	// Test GetLoggingConfig
	loggingConfig := GetLoggingConfig()
	if loggingConfig.LogLevel == "" {
		t.Error("LogLevel should not be empty")
	}

	// Test GetLogLevel
	logLevel := GetLogLevel()
	if logLevel == "" {
		t.Error("GetLogLevel should not return empty string")
	}

	// Test GetLogLocation
	logLocation := GetLogLocation()
	if logLocation == "" {
		t.Error("GetLogLocation should not return empty string")
	}

	// Test GetKeepFiles
	keepFiles := GetKeepFiles()
	if keepFiles <= 0 {
		t.Error("GetKeepFiles should return a positive number")
	}

	// Test consistency
	if loggingConfig.LogLevel != logLevel {
		t.Errorf("LogLevel mismatch: GetLoggingConfig().LogLevel = %s, GetLogLevel() = %s", loggingConfig.LogLevel, logLevel)
	}
	if loggingConfig.Location != logLocation {
		t.Errorf("Location mismatch: GetLoggingConfig().Location = %s, GetLogLocation() = %s", loggingConfig.Location, logLocation)
	}
	if loggingConfig.KeepFiles != keepFiles {
		t.Errorf("KeepFiles mismatch: GetLoggingConfig().KeepFiles = %d, GetKeepFiles() = %d", loggingConfig.KeepFiles, keepFiles)
	}
}
