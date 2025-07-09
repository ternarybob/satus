package satus

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"

	"github.com/patrickmn/go-cache"

	"github.com/gookit/color"

	"github.com/kkyr/fig"
)

// Version information - set at build time via ldflags
var (
	Version   = "dev"     // Set via -ldflags "-X main.Version=x.x.x"
	BuildTime = "unknown" // Set via -ldflags "-X main.BuildTime=..."
	GitCommit = "unknown" // Set via -ldflags "-X main.GitCommit=..."
)

type AppConfig struct {
	*serviceconfig
}

type serviceconfig struct {
	Service struct {
		Name    string `json:"name" fig:"name" default:"Go Library"`
		Version string `json:"version" fig:"version" default:"0.0.1"`
		Build   string `json:"build" fig:"build" default:"unknown"`
		Support string `json:"support,omitempty" fig:"support"`
		Slot    int    `json:"slot" fig:"slot" `
		Scope   string `json:"scope" fig:"scope" default:"DEV"`
		Port    int    `json:"port" fig:"port" default:"5001"`
		Host    string `json:"host" fig:"host" default:"localhost"`
	} `json:"service" fig:"service"`

	Logging LoggingConfig `json:"logging" fig:"logging"`

	Services []ServiceConfig `json:"services" fig:"services"`

	Connections []DataConfig `json:"connections" fig:"connections"`

	Params map[string]string `json:"-" fig:"params"`
}

type LoggingConfig struct {
	KeepFiles int    `json:"keep_files" fig:"keep_files" default:"5"`
	Location  string `json:"location" fig:"location" default:"./logs/app.log"`
	LogLevel  string `json:"loglevel" fig:"loglevel" default:"info"`
}

type ServiceConfig struct {
	Name       string `json:"name" fig:"name" default:"Service"`
	Direction  string `json:"direction" fig:"direction" default:"in"`
	Type       string `json:"type" fig:"type" default:"api"`
	Scope      string `json:"scope" fig:"scope" default:"DEV"`
	Url        string `json:"url" fig:"url" default:":8080"`
	Port       int    `json:"port" fig:"port" default:"8080"`
	Host       string `json:"host" fig:"host" default:"localhost"`
	MockUser   string `json:"mockuser" fig:"mockuser" default:"mockuser@dashs.com"`
	Connection string `json:"connection" fig:"connection" default:""`
}

type DataConfig struct {
	Name     string   `json:"name" fig:"name" default:"Connection"`
	Scope    []string `fig:"scope" default:"[default]"`
	Type     string   `fig:"type"`
	Hosts    []string `fig:"hosts" default:"[localhost]"`
	Port     string   `fig:"port"`
	Database string   `fig:"database" default:"dashs"`
	ConnStr  string   `json:"-" fig:"connstr" default:""`
	Username string   `json:"-"`
	Password string   `json:"-"`
	Cert     string   `fig:"cert"`
	Key      string   `fig:"key"`
}

const (
	ConfigName string = "config::applicationconfig"
	cacheName  string = "config::configuration"
	loglevel   string = "debug"
)

var (
	appConfig        *AppConfig     = new()
	appCache         *cache.Cache   = cache.New(5*time.Minute, 10*time.Minute)
	infoOutput                      = color.FgCyan.Render
	warnOutput                      = color.FgYellow.Render
	errorOutput                     = color.FgRed.Render
	internalloglevel zerolog.Level  = zerolog.WarnLevel
	internallog      zerolog.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger().Level(internalloglevel)
)

func new() *AppConfig {

	var cfg serviceconfig

	// Try to load config.toml first, then fall back to config.yml, or use defaults if neither exists
	configLoaded := false
	if err := fig.Load(&cfg, fig.File("config.toml")); err != nil {
		if err := fig.Load(&cfg, fig.File("config.yml")); err != nil {
			// No config file found, initialize with default values from struct tags
			// We'll manually set defaults since fig.Load() without a file doesn't work as expected
			cfg.Service.Name = "Go Library"
			cfg.Service.Version = "0.0.1"
			cfg.Service.Scope = "DEV"
			cfg.Service.Port = 5001
			cfg.Service.Host = "localhost"
			// Set default logging configuration
			cfg.Logging.LogLevel = "info"
			cfg.Logging.KeepFiles = 5
			cfg.Logging.Location = "./logs/app.log"
		} else {
			configLoaded = true
		}
	} else {
		configLoaded = true
	}

	// Log whether we loaded from file or used defaults
	if !configLoaded {
		fmt.Println(infoOutput("No config file found, using default configuration"))
	}

	// Swarm slot
	if os.Getenv("SLOT") != "" {
		cfg.Service.Slot, _ = strconv.Atoi(os.Getenv("SLOT"))
	}

	// SCope PRD|STG|DEV
	if os.Getenv("SCOPE") != "" {
		scope := strings.ToUpper(os.Getenv("SCOPE"))
		cfg.Service.Scope = scope

		for key, value := range cfg.Services {
			update := value
			value.Scope = scope

			cfg.Services[key] = update
		}
	}

	// Url of service
	if os.Getenv("URL") != "" {

		for key, value := range cfg.Services {
			update := value
			value.Url = os.Getenv("URL")

			cfg.Services[key] = update
		}

	}

	// Logging Configuration
	if os.Getenv("LOGLEVEL") != "" {
		loglevel := strings.ToUpper(os.Getenv("LOGLEVEL"))
		cfg.Logging.LogLevel = loglevel
	}

	// Apply default logging configuration if not set
	if cfg.Logging.LogLevel == "" {
		cfg.Logging.LogLevel = "info"
	}

	C := &AppConfig{&cfg}

	appCache.Set(cacheName, C, cache.DefaultExpiration)

	// fmt.Println("config:scope ->", cfg.Service.Scope)
	// fmt.Println("config:loglevel ->", cfg.Service.LogLevel)

	if _cfg, found := appCache.Get(cacheName); found {

		output := _cfg.(*AppConfig)

		return output

	}

	fmt.Println(color.Warn.Render("Cache not setup.. bup bow!"))
	return nil

}

// GetAppConfig ...
func GetAppConfig() *AppConfig {

	if _cfg, found := appCache.Get(cacheName); found {

		return _cfg.(*AppConfig)

	}

	return new()
}

// GetFirstServiceByDirection returns the first service matching the specified direction
func GetFirstServiceByDirection(direction string) (*ServiceConfig, error) {
	cfg := GetAppConfig()
	for _, service := range cfg.Services {
		if strings.EqualFold(service.Direction, direction) {
			return &service, nil
		}
	}
	return nil, fmt.Errorf("no service found with direction '%s'", direction)
}

// GetServiceDirection returns the direction for a specific service by name
func GetServiceDirection(serviceName string) string {
	cfg := GetAppConfig()
	for _, service := range cfg.Services {
		if strings.EqualFold(service.Name, serviceName) {
			if service.Direction == "" {
				return "in" // default
			}
			return service.Direction
		}
	}
	return "in" // default if service not found
}

// GetLoggingConfig returns the logging configuration
func GetLoggingConfig() LoggingConfig {
	cfg := GetAppConfig()
	return cfg.Logging
}

// GetLogLevel returns the configured log level as a string
func GetLogLevel() string {
	cfg := GetAppConfig()
	return cfg.Logging.LogLevel
}

// GetLogLocation returns the configured log file location
func GetLogLocation() string {
	cfg := GetAppConfig()
	return cfg.Logging.Location
}

// GetKeepFiles returns the number of log files to keep
func GetKeepFiles() int {
	cfg := GetAppConfig()
	return cfg.Logging.KeepFiles
}

func fileExists(path string) bool {

	_, err := os.Stat(path)
	if err == nil {
		return true
	}

	if os.IsNotExist(err) {
		return false
	}

	return false

}
