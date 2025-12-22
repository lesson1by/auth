package config

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"log"

	"authProject/internal/models"
)

func LoadConfig() (*models.Config, error) {
	var cfg models.Config
	var viperError viper.ConfigFileNotFoundError
	localViper := viper.New()
	localViper.SetConfigName("config")
	localViper.SetConfigType("yaml")
	localViper.AddConfigPath(".")
	localViper.AddConfigPath("./configs/")


	localViper.AutomaticEnv()

	localViper.SetDefault("JWT.Secret", "secret")
	localViper.SetDefault("JWT.ExpirationMinutes", 60)
	localViper.SetDefault("SERVER.Port",8080)

	if err := localViper.ReadInConfig(); err != nil {
		if errors.As(err, &viperError) {
			log.Println("Config file not found, using default settings")
		} else {
			return nil, fmt.Errorf("found config file, but encountered an error : %v", err)
		}
	}
	if err := localViper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %v", err)
	}

	if ok, err := Validate(&cfg); !ok {
		return nil, err
	}
	return &cfg, nil
}

func Validate(c *models.Config) (bool, error) {
	if c.DB.Host == "" || c.DB.Port == 0 {
		return false, errors.New("invalid database configuration")
	}
	if c.Server.Host == "" || c.Server.Port == 0 {
		return false, errors.New("invalid server configuration")
	}
	return true, nil
}
