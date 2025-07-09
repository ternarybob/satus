package main

import (
	"github.com/ternarybob/satus"
	"github.com/ternarybob/arbor"
	"github.com/ternarybob/arbor/models"
)

func main() {
	// Get satus configuration
	cfg := satus.GetAppConfig()
	
	// Initialize arbor logger with satus configuration
	logger := arbor.Logger().
		WithConsoleWriter(models.WriterConfiguration{
			Type: models.LogWriterTypeConsole,
		}).
		WithFileWriter(models.WriterConfiguration{
			Type:     models.LogWriterTypeFile,
			FileName: cfg.Logging.Location,
		}).
		WithLevelFromString(cfg.Logging.LogLevel).
		WithCorrelationId("app-startup")

	// Log some messages to demonstrate the configuration
	logger.Info().Str("service", cfg.Service.Name).Str("version", cfg.Service.Version).Msg("Application starting")
	logger.Debug().Str("log_level", cfg.Logging.LogLevel).Str("log_location", cfg.Logging.Location).Msg("Logging configuration applied")
	
	// Example of using the convenience functions
	logLevel := satus.GetLogLevel()
	logLocation := satus.GetLogLocation()
	keepFiles := satus.GetKeepFiles()
	
	logger.Info().
		Str("configured_level", logLevel).
		Str("configured_location", logLocation).
		Int("keep_files", keepFiles).
		Msg("Using satus configuration functions")
	
	// Example of updating log level at runtime
	if cfg.Service.Scope == "DEV" {
		logger.WithLevelFromString("debug").Info().Msg("Enabled debug logging for DEV environment")
	}
	
	logger.Info().Msg("Application initialization complete")
}
