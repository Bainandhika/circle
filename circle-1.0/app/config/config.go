package config

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type AppConfig struct {
	Host          string `yaml:"host"`
	Port          int    `yaml:"port"`
	SecretKey     string `yaml:"secret-key"`
	LogPath       string `yaml:"log-path"`
	BillPath      string `yaml:"bill-path"`
	FontAssetPath string `yaml:"font-asset-path"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type config struct {
	App      AppConfig      `yaml:"app"`
	Database DatabaseConfig `yaml:"database"`
	Redis    RedisConfig    `yaml:"redis"`
}

var (
	App   AppConfig
	DB    DatabaseConfig
	Redis RedisConfig
)

const (
	configPath = "circle-config.yaml"
)

func loadConfig() error {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("error read %s, err: %v", configPath, err)
	}

	var config config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return fmt.Errorf("error unmarshal config data, err: %v", err)
	}

	DB = config.Database
	App = config.App
	Redis = config.Redis

	return nil
}

func init() {
	err := loadConfig()
	if err != nil {
		log.Println("Error load config, err: " + err.Error())
		panic(err)
	}
}
