package config

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string     `yaml:"env"`
	StoragePath string     `yaml:"storage_path"`
	HTTPServer  HTTPServer `yaml:"http_server"`
	Database    Database   `yaml:"database"`
}

type HTTPServer struct {
	Address     string        `yaml:"address"`
	Timeout     time.Duration `yaml:"timeout"`
	IdleTimeout time.Duration `yaml:"idle_timeout"`
	User        string        `yaml:"user"`
	Password    string        `yaml:"password"`
}

type Database struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Dbname   string `yaml:"dbname"`
}

func MustLoad() *Config {
	configPath := filepath.Join("./config/config.yaml")

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file %s does not exist", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("error reading config file: %s", err)
	}

	log.Printf("Loaded config: %+v", cfg)

	return &cfg
}
