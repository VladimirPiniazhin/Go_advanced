package configs

import (
	"go/order-api/pkg/jwt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MailConf EmailConfig
	Db       DbConfig
	Jwt      jwt.JWT
}

type DbConfig struct {
	Dsn string
}

type EmailConfig struct {
	Email    string
	Password string
	Address  string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, using default config")
	}
	return &Config{
		MailConf: EmailConfig{
			Email:    os.Getenv("EMAIL"),
			Password: os.Getenv("PASSWORD"),
			Address:  os.Getenv("ADDRESS"),
		},
		Db: DbConfig{
			Dsn: os.Getenv("DB_DSN"),
		},
		Jwt: jwt.JWT{
			Secret: os.Getenv("SECRET"),
		},
	}
}
