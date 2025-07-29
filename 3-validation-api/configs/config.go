package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	mail MailConfig
}

type MailConfig struct {
	email    string
	password string
	address  string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, using default config")
	}
	return &Config{
		mail: MailConfig{
			email:    os.Getenv("EMAIL"),
			password: os.Getenv("PASSWORD"),
		},
	}
}
