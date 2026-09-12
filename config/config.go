package config

import (
	"log"

	"github.com/spf13/viper"
)

var CLOUDINARY_BASE_URL string

type Config struct {
	App        App
	JWTConfig  JWTConfig
	Postgres   Postgres
	SMTP       SMTP
	Cloudinary Cloudinary
}

type App struct {
	Name        string
	Environment string
	Host        int
	Port        int
	URL         string
}

type JWTConfig struct {
	Secret string
}

type Postgres struct {
	Host     string
	Port     string
	Database string
	User     string
	Password string
	SSLMode  string
	URI      string
}

type SMTP struct {
	User      string
	Pass      string
	Host      string
	Port      int
	EmailFrom string
}

type Cloudinary struct {
	CloudName string
	APIKey    string
	APISecret string
	BaseURL   string
}

type Google struct {
	ClientID string
}

func New() *Config {
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("[config-file-fail-load] \n", err.Error())
	}

	v := viper.GetViper()
	viper.AutomaticEnv()

	CLOUDINARY_BASE_URL = v.GetString("CLOUDINARY_BASE_URL")

	return &Config{
		App: App{
			Name:        v.GetString("APP_NAME"),
			Environment: v.GetString("APP_ENV"),
			Host:        v.GetInt("APP_HOST"),
			Port:        v.GetInt("APP_PORT"),
			URL:         v.GetString("APP_URL"),
		},
		JWTConfig: JWTConfig{
			Secret: v.GetString("JWT_SECRET"),
		},
		Postgres: Postgres{
			Host:     v.GetString("POSTGRES_HOST"),
			Port:     v.GetString("POSTGRES_PORT"),
			Database: v.GetString("POSTGRES_DATABASE"),
			User:     v.GetString("POSTGRES_USER"),
			Password: v.GetString("POSTGRES_PASS"),
			SSLMode:  v.GetString("POSTGRES_SSL_MODE"),
			URI:      v.GetString("POSTGRES_URI"),
		},
		SMTP: SMTP{
			User:      v.GetString("SMTP_USER"),
			Pass:      v.GetString("SMTP_PASSWORD"),
			Host:      v.GetString("SMTP_HOST"),
			Port:      v.GetInt("SMTP_PORT"),
			EmailFrom: v.GetString("EMAIL_FROM"),
		},
		Cloudinary: Cloudinary{
			CloudName: v.GetString("CLOUDINARY_NAME"),
			APIKey:    v.GetString("CLOUDINARY_API_KEY"),
			APISecret: v.GetString("CLOUDINARY_API_SECRET"),
			BaseURL:   CLOUDINARY_BASE_URL,
		},
	}
}
