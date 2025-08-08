package main

import (
	"fmt"
	"go/order-api/configs"
	"go/order-api/internals/auth"
	"go/order-api/internals/link"
	"go/order-api/internals/order"
	"go/order-api/internals/product"
	"go/order-api/internals/stat"
	"go/order-api/internals/user"
	"go/order-api/internals/verify"
	"go/order-api/pkg/db"
	"go/order-api/pkg/email"
	"go/order-api/pkg/event"
	"go/order-api/pkg/jwt"
	"go/order-api/pkg/middleware"
	"net/http"
)

// App создает и настраивает HTTP-хендлер для тестирования
func App() http.Handler {
	// Загружаем конфигурацию
	config := configs.LoadConfig()

	// Создаём подключение к БД
	database := db.NewDb(config)

	// Создаём роутер
	router := http.NewServeMux()

	// Создаём EventBus
	eventBus := event.NewEventBus()

	// Создаём сервисы
	jwtService := jwt.NewJWT(config.Jwt.Secret)
	emailService := email.NewEmailService(config.MailConf.Email, config.MailConf.Password, config.MailConf.Address)

	// Создаём репозитории
	linkRepository := link.NewLinkRepository(database)
	productRepository := product.NewProductRepository(database)
	userRepository := user.NewUserRepository(database)
	statRepository := stat.NewStatRepository(database)
	orderRepository := order.NewOrderRepository(database)

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
	order.NewOrderHandler(router, order.OrderHandlerDeps{
		OrderRepository: orderRepository,
		Config:          config,
		UserRepository:  userRepository,
	})

	// Создаём middleware stack только с CORS и логированием
	stack := middleware.Chain(
		middleware.Simple(middleware.CORSSimple),
		middleware.Simple(middleware.LoggingSimple),
	)

	// Запускаем статистический сервис в фоне
	go statService.AddClick()

	return stack(router, config)
}

func main() {
	server := http.Server{
		Addr:    ":8081",
		Handler: App(),
	}

	fmt.Println("Server listening on port 8081")
	server.ListenAndServe()
}
