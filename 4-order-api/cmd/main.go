package main

import (
	"fmt"
	"go/order-api/configs"
	"go/order-api/internals/auth"
	"go/order-api/internals/link"
	"go/order-api/internals/product"
	"go/order-api/internals/stat"
	"go/order-api/internals/user"
	"go/order-api/internals/verify"
	"go/order-api/pkg/db"
	"go/order-api/pkg/email"
	"go/order-api/pkg/event"
	"go/order-api/pkg/jwt"
	"go/order-api/pkg/middleware"
	"log"
	"net/http"
	"os"

	logger "github.com/sirupsen/logrus"
)

func main() {
	// Загружаем конфигурацию
	config := configs.LoadConfig()

	// Создаём подключение к БД
	database := db.NewDb(config)

	// Создаём роутер
	router := http.NewServeMux()

	// Создаём EventBus
	eventBus := event.NewEventBus()

	// Logging setup
	file, err := os.OpenFile("logs.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Не удалось открыть файл логов: %v", err)
	}
	logger.SetOutput(file)
	logger.SetFormatter(&logger.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	logger.SetLevel(logger.InfoLevel)

	// Создаём сервисы
	jwtService := jwt.NewJWT(config.Jwt.Secret)
	emailService := email.NewEmailService(config.MailConf.Email, config.MailConf.Password, config.MailConf.Address)

	// Создаём репозитории
	linkRepository := link.NewLinkRepository(database)
	productRepository := product.NewProductRepository(database)
	userRepository := user.NewUserRepository(database)
	statRepository := stat.NewStatRepository(database)

	// Создаём бизнес-сервисы
	authService := auth.NewAuthService(userRepository, jwtService)
	statService := stat.NewStatService(&stat.StatServiceDeps{
		EventBus:       eventBus,
		StatRepository: statRepository,
	})

	// Создаём хендлеры
	stat.NewStatHandler(router, stat.StatHandlerDeps{
		StatRepository: statRepository,
	})
	auth.NewAuthHandler(router, auth.AuthHandlerDeps{
		AuthService: authService,
	})
	link.NewLinkHandler(router, link.LinkHandlerDeps{
		LinkRepository: linkRepository,
		EventBus:       eventBus,
		Config:         config,
	})
	verify.NewVerifyHandler(router, verify.VerifyHandlerDeps{
		EmailService: emailService,
		Config:       config,
	})
	product.NewProductHandler(router, product.ProductHandlerDeps{
		ProductRepository: productRepository,
		Config:            config,
	})

	// Создаём middleware stack только с CORS и логированием
	stack := middleware.Chain(
		middleware.Simple(middleware.CORSSimple),
		middleware.Simple(middleware.LoggingSimple),
	)

	server := http.Server{
		Addr:    ":8081",
		Handler: stack(router, config),
	}

	// Запускаем статистический сервис
	go statService.AddClick()

	fmt.Println("Server listening on port 8081")
	server.ListenAndServe()

}
