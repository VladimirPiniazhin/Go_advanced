package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Db       DbConfig
	AuthConf AuthConfig
}

type DbConfig struct {
	Db string
}

type AuthConfig struct {
	Secret string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, using default config")
	}
	return &Config{
		Db: DbConfig{
			Db: os.Getenv("DSN"),
		},
		AuthConf: AuthConfig{
			Secret: os.Getenv("TOKEN"),
		},
	}
}
