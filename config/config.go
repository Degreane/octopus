package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the core configuration settings for the application.
// It contains sensitive and environment-specific parameters including database connection details,
// authentication secrets, and runtime environment specification.
type Config struct {
	MongoURI    string
	JWTSecret   string
	SessionKey  string
	Environment string
}

// Route defines the configuration for a single route in the application.
// It specifies the HTTP method, path, optional pre-check validations, and associated view.
type Route struct {
	Method    string  `yaml:"method"`
	Path      string  `yaml:"path"`
	PreCheck  []Check `yaml:"preCheck"`
	View      string  `yaml:"view"`
	WebSocket bool    `yaml:"websocket,omitempty"`
}

// Check represents a pre-check configuration for routes, containing headers and script validation details.
// It is used to define custom validation or preprocessing steps before executing a route handler.
type Check struct {
	Headers string `yaml:"headers"`
	Script  string `yaml:"script"`
}

// ModulesConfig represents the configuration for a module in the application.
// It contains metadata and configuration details for individual modules, including
// their identification, dependencies, routes, and other settings.
type ModulesConfig struct {
	Name         string   `yaml:"Name"`
	Description  string   `yaml:"Description"`
	Version      string   `yaml:"Version"`
	Author       string   `yaml:"Author"`
	Email        string   `yaml:"Email"`
	Website      string   `yaml:"Website"`
	License      string   `yaml:"License"`
	Dependencies []string `yaml:"Dependencies"`
	Settings     []string `yaml:"Settings"`
	BasePath     string   `yaml:"BasePath"`
	LocalPath    string   `yaml:"LocalPath,omitempty"`
	DB           string   `yaml:"db"`
	Routes       []Route  `yaml:"Routes,omitempty"`
}

type Storage string

const (
	MongoDB    Storage = "mongodb"
	Redis      Storage = "redis"
	Memory     Storage = "memory"
	File       Storage = "file"
	MySql      Storage = "mysql"
	Postgresql Storage = "postgresql"
)

// AppConfig represents the server configuration settings for the application.
// It defines key parameters for server initialization and runtime behavior.
type AppConfig struct {
	Port         string  `yaml:"Port"`
	Prefork      bool    `yaml:"Prefork"`
	Storage      Storage `yaml:"Storage"`
	Debug        bool    `yaml:"Debug"`
	ServerHeader string  `yaml:"ServerHeader"`
}

// New creates and returns a new Config instance with environment-specific configuration values.
// It retrieves configuration values from environment variables for MongoDB URI, JWT secret,
// session key, and environment type.
func New() *Config {
	return &Config{
		MongoURI:    os.Getenv("MONGO_URI"),
		JWTSecret:   os.Getenv("JWT_SECRET"),
		SessionKey:  os.Getenv("SESSION_KEY"),
		Environment: os.Getenv("ENV"),
	}
}

// ParseModulesConfig reads and parses the modules configuration from the modules.yaml file.
// It returns a slice of ModulesConfig and an error if the file cannot be read or parsed.
func ParseModulesConfig() ([]ModulesConfig, error) {
	// Read the YAML file
	// Read the modules configuration file from the specified path
	data, err := os.ReadFile("config/modules.yaml")
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}

	// Parse YAML into ModulesConfig slice
	// modules stores the parsed configuration for multiple modules from the modules configuration file
	var modules []ModulesConfig
	err = yaml.Unmarshal(data, &modules)
	if err != nil {
		return nil, fmt.Errorf("error parsing config file: %v", err)
	}

	return modules, nil
}

// ParseServerConfig reads and parses the server configuration from the config.yaml file.
// It returns a pointer to an AppConfig struct and an error if the file cannot be read or parsed.
func ParseServerConfig() (*AppConfig, error) {
	// Read the YAML file
	data, err := os.ReadFile("config/config.yaml")
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}
	var appConfig *AppConfig
	err = yaml.Unmarshal(data, &appConfig)
	if err != nil {
		return nil, fmt.Errorf("error parsing config file: %v", err)
	}
	return appConfig, nil
}
