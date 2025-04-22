package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	MongoURI    string
	JWTSecret   string
	SessionKey  string
	Environment string
}

type Route struct {
	Method   string  `yaml:"method"`
	Path     string  `yaml:"path"`
	PreCheck []Check `yaml:"preCheck"`
	View     string  `yaml:"view"`
}

type Check struct {
	Headers string `yaml:"headers"`
	Script  string `yaml:"script"`
}

type ModulesConfig struct {
	Name     string  `yaml:"Name"`
	BasePath string  `yaml:"BasePath"`
	DB       string  `yaml:"db"`
	Routes   []Route `yaml:"Routes,omitempty"`
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

type AppConfig struct {
	Port    string  `yaml:"Port"`
	Prefork bool    `yaml:"Prefork"`
	Storage Storage `yaml:"Storage"`
	Debug   bool    `yaml:"Debug"`
}

func New() *Config {
	return &Config{
		MongoURI:    os.Getenv("MONGO_URI"),
		JWTSecret:   os.Getenv("JWT_SECRET"),
		SessionKey:  os.Getenv("SESSION_KEY"),
		Environment: os.Getenv("ENV"),
	}
}

func ParseModulesConfig() ([]ModulesConfig, error) {
	// Read the YAML file
	data, err := os.ReadFile("config/modules.yaml")
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}

	// Parse YAML into ModulesConfig slice
	var modules []ModulesConfig
	err = yaml.Unmarshal(data, &modules)
	if err != nil {
		return nil, fmt.Errorf("error parsing config file: %v", err)
	}

	return modules, nil
}

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
