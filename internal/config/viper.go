package config

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	JWTSecret  string `mapstructure:"JWT_SECRET"`
	Port       string `mapstructure:"PORT"`
	Env        string `mapstructure:"ENV"`
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`
}

func LoadConfig() (cfg Config, err error) {
	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		log.Printf("Warning: couldn't find .env file, reading from OS environment: %v", err)
	}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		return cfg, err
	}

	if strings.TrimSpace(cfg.JWTSecret) == "" {
		return cfg, fmt.Errorf("JWT_SECRET must be set")
	}
	if len(cfg.JWTSecret) < 32 {
		return cfg, fmt.Errorf("JWT_SECRET must be at least 32 characters")
	}

	return cfg, nil
}
