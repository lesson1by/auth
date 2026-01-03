package models

type Config struct {
	JWT struct {
		Secret            string `mapstructure:"secret"`
		ExpirationMinutes int    `mapstructure:"expiration_minutes"`
	} `mapstructure:"jwt"`
	Server struct {
		Port int    `mapstructure:"port"`
		Host string `mapstructure:"host"`
	} `mapstructure:"server"`
	DB struct {
		Name     string `mapstructure:"name"`
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
	} `mapstructure:"db"`
}
