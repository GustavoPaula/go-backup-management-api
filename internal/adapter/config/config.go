package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DB    *DB
	HTTP  *HTTP
	Token *Token
}

type DB struct {
	User string
	Pass string
	Host string
	Port string
	Name string
}

type HTTP struct {
	Port string
}

type Token struct {
	Duration     string
	JwtSecretKey string
}

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

	token := &Token{
		Duration:     os.Getenv("TOKEN_DURATION"),
		JwtSecretKey: os.Getenv("JWT_SECRET_KEY"),
	}

	return &Config{
		db,
		http,
		token,
	}, nil
}
