package config

import (
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
	DefaultConfigPath = "circle-config.env"
)

func InitConfig(path string) {
	if path == "" {
		path = DefaultConfigPath
	}

	err := godotenv.Load(path)
	if err != nil {
		log.Fatalln("Error loading .env file: " + err.Error())
	}

	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalln("Error parse env data to struct: " + err.Error())
	}

	App = cfg.App
	DB = cfg.DB
	Redis = cfg.Redis
}
