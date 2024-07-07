package config

import (
	"errors"
	"fmt"
	"log"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type AppConfig struct {
	Host          string `env:"APP_HOST"`
	Port          int    `env:"APP_PORT"`
	SecretKey     string `env:"APP_SECRET_KEY"`
	LogPath       string `env:"APP_LOG_PATH"`
	BillPath      string `env:"APP_BILL_PATH"`
	FontAssetPath string `env:"APP_FONT_ASSET_PATH"`
}

type DatabaseConfig struct {
	Host     string `env:"DATABASE_HOST"`
	Port     int    `env:"DATABASE_PORT"`
	Username string `env:"DATABASE_USERNAME"`
	Password string `env:"DATABASE_PASSWORD"`
	Name     string `env:"DATABASE_NAME"`
}

type RedisConfig struct {
	Host     string `env:"REDIS_HOST"`
	Port     int    `env:"REDIS_PORT"`
	Username string `env:"REDIS_USERNAME"`
	Password string `env:"REDIS_PASSWORD"`
}

type config struct {
	App   AppConfig
	DB    DatabaseConfig
	Redis RedisConfig
}

var (
	App   AppConfig
	DB    DatabaseConfig
	Redis RedisConfig
)

const (
	configPath = "circle-config.env"
)

func loadConfig() error {
	err := godotenv.Load(configPath)
	if err != nil {
		return errors.New("error loading .env file")
	}

	// Parse the environment variables into the struct
	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		return fmt.Errorf("failed to parse env into struct: %v", err)
	}

	App = cfg.App
	DB = cfg.DB
	Redis = cfg.Redis

	// Print the config struct to verify the values
	fmt.Printf("Config: %+v\n", cfg)

	return nil
}

func init() {
	err := loadConfig()
	if err != nil {
		log.Println("Error load config, err: " + err.Error())
		panic(err)
	}
}
