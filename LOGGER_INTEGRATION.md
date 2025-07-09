# Logger Configuration Integration

This document summarizes the integration of satus configuration with arbor logger configuration changes.

## Changes Made

### Satus Configuration Updates

1. **Added `LoggingConfig` struct** with fields:
   - `KeepFiles`: Number of log files to keep
   - `Location`: Log file location
   - `LogLevel`: Log level string

2. **Updated `serviceconfig` struct** to include:
   - `Logging LoggingConfig` field
   - Removed `LogLevel` and `LogFile` from service section

3. **Updated `ServiceConfig` struct** to remove:
   - `LogLevel` field (moved to logging section)
   - `LogFile` field (moved to logging section)

4. **Added convenience functions**:
   - `GetLoggingConfig()` - Returns full logging configuration
   - `GetLogLevel()` - Returns log level string
   - `GetLogLocation()` - Returns log file location
   - `GetKeepFiles()` - Returns number of files to keep

5. **Updated configuration loading**:
   - Environment variable `LOGLEVEL` now sets `cfg.Logging.LogLevel`
   - Default logging configuration is set when no config file exists

### Arbor Integration

The satus configuration now integrates with arbor's `WithLevelFromString()` method:

```go
// Get satus configuration
cfg := satus.GetAppConfig()

// Apply to arbor logger
logger := arbor.Logger().
    WithConsoleWriter(models.WriterConfiguration{
        Type: models.LogWriterTypeConsole,
    }).
    WithLevelFromString(cfg.Logging.LogLevel)
```

### Configuration File Structure

The configuration file structure has been updated:

#### Before (Old Structure)
```toml
[service]
name = "my-service"
loglevel = "info"
logfile = "./logs/app.log"
```

#### After (New Structure)
```toml
[service]
name = "my-service"
rest_loglevel = "info"  # Only for REST API logging

[logging]
keep_files = 5
location = "./logs/app.log"
loglevel = "info"
```

### Environment Variables

- `LOGLEVEL`: Sets the logging level for the application
- Other environment variables unchanged

### Omnis Integration

Updated omnis to use the new configuration structure:

```go
// Added new function to get arbor logger with satus config
func getArborLogger() arbor.ILogger {
    return arbor.Logger().
        WithConsoleWriter(models.WriterConfiguration{
            Type: models.LogWriterTypeConsole,
        }).
        WithLevelFromString(satus.GetLogLevel()).
        WithPrefix("omnis")
}
```

## Benefits

1. **Separation of Concerns**: Arbor is now purely a logging library without application-specific configuration
2. **Flexibility**: Applications can handle their own configuration and pass strings to arbor
3. **Order Independence**: Logger starts with INFO level and can be updated when config loads
4. **Consistency**: All projects use the same logging configuration pattern

## Migration Guide

For applications using satus configuration:

1. Update configuration files to use the new `[logging]` section
2. Remove `loglevel` and `logfile` from service configurations
3. Use `satus.GetLogLevel()` to get the configured log level
4. Apply to arbor logger using `WithLevelFromString()`

Example:
```go
// Old approach (no longer available)
// logger := arbor.Logger().ApplyConfig(cfg)

// New approach
logger := arbor.Logger().
    WithConsoleWriter(config).
    WithLevelFromString(satus.GetLogLevel())
```

## Testing

All tests pass for:
- Satus configuration loading and access functions
- Arbor logger integration
- Omnis project integration
- Backward compatibility with existing functionality
