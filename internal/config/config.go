package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env       string    `yaml:"env" env-default:"local"`
	Servers   Servers   `yaml:"servers"`
	PrimaryDB PrimaryDB `yaml:"database"`
}
type Servers struct {
	Grpc GRPCServer `yaml:"grpc"`
}
type GRPCServer struct {
	Port    int           `yaml:"port" `
	Timeout time.Duration `yaml:"timeout" env-default:"60s"`
}
type PrimaryDB struct {
	Host        string `yaml:"host"`
	Port        string `yaml:"port"`
	User        string `yaml:"user"`
	Dbname      string `yaml:"dbname"`
	MaxAttempts int    `yaml:"max_attempts"`
}

func MustLoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set in environment variables")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Config file does not exist: %s", configPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	return &cfg
}
