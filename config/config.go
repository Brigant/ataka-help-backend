package config

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Server struct {
	AppAddress      string
	AppReadTimeout  time.Duration
	AppWriteTimeout time.Duration
	AppIdleTimeout  time.Duration
}

type Postgres struct {
	Host     string
	Port     string
	Database string
	User     string
	Password string
	SSLmode  string
}

type Config struct {
	Server          Server
	Postgres        Postgres
	LogLevel        string
	Salt            string
	SigningKey      string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

// Allowed logger levels & config key.
const (
	DebugLogLvl = "DEBUG"
	InfoLogLvl  = "INFO"
	ErrorLogLvl = "ERROR"
)

var errNotAllowedLoggelLevel = errors.New("not allowed logger level")

func InitConfig() (Config, error) {
	// Set the path of the .env file
	viper.SetConfigFile(".env")

	// Add the path from which to try and retrieve the .env file
	viper.AddConfigPath("../.")

	// Enable reading variables from the .env file
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return Config{}, fmt.Errorf("error while reading config: %w", err)
	}

	logLevel := viper.GetString("LOG_LEVEL")

	if err := validate(logLevel); err != nil {
		return Config{}, fmt.Errorf("error while init config: %w", err)
	}

	config := Config{
		Server: Server{
			AppAddress:      ":" + viper.GetString("APP_PORT"),
			AppReadTimeout:  viper.GetDuration("APP_READ_TIMEOUT") * time.Second,
			AppWriteTimeout: viper.GetDuration("APP_WRITE_TIMEOUT") * time.Second,
			AppIdleTimeout:  viper.GetDuration("APP_IDLE_TIMEOUT") * time.Second,
		},
		LogLevel: logLevel,
		Postgres: Postgres{
			Host:     viper.GetString("DB_HOST"),
			Port:     viper.GetString("DB_PORT"),
			Database: viper.GetString("DB_NAME"),
			User:     viper.GetString("DB_USER"),
			Password: viper.GetString("DB_PASSWORD"),
			SSLmode:  viper.GetString("DB_SSL_MODE"),
		},
		Salt:            viper.GetString("APP_SALT"),
		SigningKey:      viper.GetString("SIGNING_KEY"),
		AccessTokenTTL:  viper.GetDuration("ACCESS_TOKEN_TTL") * time.Minute,
		RefreshTokenTTL: viper.GetDuration("REFRESH_TOKEN_TTL") * time.Hour,
	}

	return config, nil
}

func validate(logLevel string) error {
	if strings.ToUpper(logLevel) != DebugLogLvl &&
		strings.ToUpper(logLevel) != ErrorLogLvl &&
		strings.ToUpper(logLevel) != InfoLogLvl {
		return errNotAllowedLoggelLevel
	}

	return nil
}
