package config

import (
	"os"

	"github.com/joho/godotenv"
)

type (
	Config struct {
		DB   *DB
		HTTP *HTTP
	}

	DB struct {
		User string
		Pass string
		Host string
		Port string
		Name string
	}

	HTTP struct {
		Port string
	}
)

func New() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	db := &DB{
		User: os.Getenv("DB_USER"),
		Pass: os.Getenv("DB_PASS"),
		Host: os.Getenv("DB_HOST"),
		Port: os.Getenv("DB_PORT"),
		Name: os.Getenv("DB_NAME"),
	}

	http := &HTTP{
		Port: os.Getenv("HTTP_PORT"),
	}

	return &Config{
		db,
		http,
	}, nil
}
