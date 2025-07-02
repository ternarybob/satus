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

type AppConfig struct {
	*serviceconfig
}

type serviceconfig struct {
	Service struct {
		Name     string `json:"name" fig:"name" default:"Go Library"`
		Version  string `json:"version" fig:"version" default:"0.0.1"`
		Support  string `json:"support,omitempty" fig:"support"`
		Slot     int    `json:"slot" fig:"slot" `
		Scope    string `json:"scope" fig:"scope" default:"DEV"`
		Port     int    `json:"port" fig:"port" default:"5001"`
		Host     string `json:"host" fig:"host" default:"localhost"`
		LogLevel     string `json:"loglevel" fig:"loglevel" default:"info"`
		RestLogLevel string `json:"rest_loglevel" fig:"rest_loglevel" default:"info"`
		LogFile      string `json:"logfile" fig:"logfile" default:""`
	} `json:"service" fig:"service"`

	Services []ServiceConfig `json:"services" fig:"services"`

	Connections []DataConfig `json:"connections" fig:"connections"`

	Params map[string]string `json:"-" fig:"params"`
}

type ServiceConfig struct {
	Name       string `json:"name" fig:"name" default:"Service"`
	Direction  string `json:"direction" fig:"direction" default:"in"`
	Type       string `json:"type" fig:"type" default:"api"`
	Scope      string `json:"scope" fig:"scope" default:"DEV"`
	Url        string `json:"url" fig:"url" default:":8080"`
	Port       int    `json:"port" fig:"port" default:"8080"`
	Host       string `json:"host" fig:"host" default:"localhost"`
	LogLevel   string `json:"loglevel" fig:"loglevel" default:"info"`
	LogFile    string `json:"logfile" fig:"logfile" default:""`
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

	// Try to load config.toml first, then fall back to config.yml
	if err := fig.Load(&cfg, fig.File("config.toml")); err != nil {
		if err := fig.Load(&cfg, fig.File("config.yml")); err != nil {
			panic(fmt.Sprintf("config file not found: tried config.toml and config.yml - %v", err))
		}
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

	// Logging Level
	if os.Getenv("LOGLEVEL") != "" {
		loglevel := strings.ToUpper(os.Getenv("LOGLEVEL"))
		cfg.Service.LogLevel = loglevel

		for key, value := range cfg.Services {
			update := value
			value.LogLevel = loglevel

			cfg.Services[key] = update
		}
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
