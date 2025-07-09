# Configuration

Satus supports configuration through TOML and YAML files. The configuration allows you to set log levels and other settings that affect the behavior of the application.

## Configuration File Format

The configuration file supports both TOML (`config.toml`) and YAML (`config.yml`) formats. TOML is preferred and will be loaded first if both exist.

### TOML Format Example (`config.toml`):

```toml
[service]
name = "my-service"
version = "1.0.0"
support = "support@example.com"
slot = 1
scope = "DEV"
port = 5001
host = "localhost"
rest_loglevel = "info"

[logging]
keep_files = 5
location = "./logs/app.log"
loglevel = "info"

[[services]]
name = "test-service"
direction = "in"
type = "api"
scope = "DEV"
url = ":8080"
mockuser = "test@example.com"
connection = ""

[[connections]]
name = "default"
scope = ["DEV"]
type = "sqlite"
hosts = ["localhost"]
port = "5432"
database = "test"
connstr = ""
cert = ""
key = ""

[params]
test_param = "test_value"
```

## Configuration Sections

### Service Configuration

The `[service]` section contains general application configuration:

- `name`: Application name
- `version`: Application version
- `support`: Support contact information
- `slot`: Service slot number (for swarm deployments)
- `scope`: Environment scope (DEV, STG, PRD)
- `port`: HTTP port to bind to
- `host`: Host to bind to
- `rest_loglevel`: Log level for REST API logging

### Logging Configuration

The `[logging]` section controls application logging:

- `keep_files`: Number of log files to keep (integer)
- `location`: Path to the log file (string)
- `loglevel`: Log level to use (string)

### Supported Log Levels

- `"trace"` - Most verbose, includes all log messages
- `"debug"` - Debug messages and above
- `"info"` - Informational messages and above (default)
- `"warn"` or `"warning"` - Warning messages and above
- `"error"` - Error messages and above
- `"fatal"` - Fatal messages and above
- `"panic"` - Only panic messages
- `"disabled"` or `"off"` - No logging

## Usage

### Basic Usage

```go
import (
    "github.com/ternarybob/satus"
    "github.com/ternarybob/arbor"
    "github.com/ternarybob/arbor/models"
)

func main() {
    // Get configuration
    cfg := satus.GetAppConfig()
    
    // Apply to arbor logger
    logger := arbor.Logger().
        WithConsoleWriter(models.WriterConfiguration{
            Type: models.LogWriterTypeConsole,
        }).
        WithLevelFromString(cfg.Logging.LogLevel)
    
    logger.Info().Msg("Application started")
}
```

### Accessing Configuration

```go
// Get full configuration
cfg := satus.GetAppConfig()

// Get specific logging configuration
loggingConfig := satus.GetLoggingConfig()

// Get specific values
logLevel := satus.GetLogLevel()
logLocation := satus.GetLogLocation()
keepFiles := satus.GetKeepFiles()
```

## Environment Variables

Certain configuration values can be overridden with environment variables:

- `SLOT`: Override service slot number
- `SCOPE`: Override environment scope
- `URL`: Override service URL
- `LOGLEVEL`: Override log level

## Order Independence

The configuration system is designed to be order-independent. Configuration is loaded once and cached, so you can access it from anywhere in your application without worrying about load order.
