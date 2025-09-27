package gojs

import (
	"log"
	"os"

	"github.com/goccy/go-yaml"
)

type ServerOptions struct {
	Server    bool
	EnableWeb bool
	WebDir    string
}

type Options struct {
	ServerOptions
}

type GoJSConfig struct {
	Title    string `yaml:"title"`
	Server   bool   `yaml:"server"`
	Port     int    `yaml:"port"`
	Endpoint string `yaml:"endpoint"`
	RedisKey string `yaml:"redisKey"`
	WebDir   string `yaml:"webDir"` // web page static directory
}

type Config struct {
	Initialized bool
	MySQL       MySQLConfig `yaml:"mysql"` // Map of MySQL configs by name
	Redis       RedisConfig `yaml:"redis"` // Redis config
	GOJS        GoJSConfig  `yaml:"gojs"`  // Gojs config
}

// MySQLConfig holds MySQL database connection details
type MySQLConfig struct {
	Name       string `yaml:"name,omitempty"`
	ConnString string `yaml:"connString"`
	Host       string `yaml:"host"`
	Port       int    `yaml:"port"`
	User       string `yaml:"user"`
	Password   string `yaml:"password"`
	DB         string `yaml:"db"`
	Timeout    int    `yaml:"timeout"`
}

type RedisConfig struct {
	Addr     string `yaml:"addr"`
	DB       int    `yaml:"db"`
	DBConfig int    `yaml:"dbConfig"`
	Password string `yaml:"password,omitempty"`
}

// LoadConfig loads configuration from a file
func LoadConfig(yamlFilePath string) *Config {
	cfg := Config{}

	data, err := os.ReadFile(yamlFilePath)
	if err != nil {
		log.Printf("[gojs] %s: %v", "Warning: Could not read YAML config file", err)
		return nil
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Printf("[gojs] %s: %v", "Warning: Could not parse YAML config file", err)
		return nil
	}

	return &cfg
}
