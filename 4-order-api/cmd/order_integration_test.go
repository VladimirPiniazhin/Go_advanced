package main

import (
	"bytes"
	"encoding/json"
	"go/order-api/internals/auth"
	"go/order-api/internals/order"
	"go/order-api/internals/product"
	"go/order-api/internals/user"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var testUserEmail = "test_order_user@mail.ru"
var testUserPassword = "111111"

func initTestDb() *gorm.DB {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(postgres.Open(os.Getenv("DB_DSN")), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}

func initTestData(db *gorm.DB) (uint, []uint) {
	// Тестовый пользователь
	testUser := &user.User{
		Email:    testUserEmail,
		Password: "$2a$10$8cnP5MJNu/OUyTEZhZlxaOl9qE2ohfRHWMJYcho8OKuTBp7ZsqeBa", // пароль: 111111
		Name:     "Test Order User",
		Phone:    "+79123456789",
	}
	db.Create(testUser)

	// Тестовые продукты
	products := []product.Product{
		{
			Name:        "Test Product 1",
			Description: "Test Description 1",
			Price:       1000,
		},
		{
			Name:        "Test Product 2",
			Description: "Test Description 2",
			Price:       2000,
		},
		{
			Name:        "Test Product 3",
			Description: "Test Description 3",
			Price:       1500,
		},
	}

	var productIDs []uint
	for _, p := range products {
		db.Create(&p)
		productIDs = append(productIDs, p.ID)
	}

	return testUser.ID, productIDs
}

func removeTestData(db *gorm.DB) {
	// Удаляем тестовые данные
	db.Unscoped().Where("email = ?", testUserEmail).Delete(&user.User{})
	db.Unscoped().Where("name LIKE ?", "Test Product%").Delete(&product.Product{})
	// Удаляем все заказы тестового пользователя
	db.Unscoped().Where("user_id IN (SELECT id FROM users WHERE email = ?)", testUserEmail).Delete(&order.Order{})
}

func getAuthToken(t *testing.T, testServer *httptest.Server) string {
	// Логинимся и получаем токен
	loginData, _ := json.Marshal(&auth.LoginRequest{
		Email:    testUserEmail,
		Password: testUserPassword,
	})

	response, err := http.Post(testServer.URL+"/auth/login", "application/json", bytes.NewReader(loginData))
	if err != nil {
		t.Fatal(err)
	}

	if response.StatusCode != 200 {
		t.Fatalf("Login failed with status %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}

	var authResponse auth.AuthorizationResponse
	err = json.Unmarshal(body, &authResponse)
	if err != nil {
		t.Fatal(err)
	}

	return authResponse.Token
}

func TestOrderIntegration(t *testing.T) {
	// Инициализация БД
	db := initTestDb()
	userID, productIDs := initTestData(db)
	testServer := httptest.NewServer(App())
	defer testServer.Close()

	// Получаем токен авторизации
	token := getAuthToken(t, testServer)

	// Тест 1: Создание заказа
	t.Run("CreateOrder", func(t *testing.T) {
		orderData := order.OrderCreateRequest{
			OrderItems: []order.OrderItem{
				{
					ProductID: productIDs[0],
					Quantity:  2,
				},
				{
					ProductID: productIDs[1],
					Quantity:  1,
				},
			},
		}

		data, _ := json.Marshal(orderData)

		req, _ := http.NewRequest("POST", testServer.URL+"/order", bytes.NewReader(data))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		client := &http.Client{}
		response, err := client.Do(req)

		if err != nil {
			t.Fatal(err)
		}

		if response.StatusCode != 201 {
			body, _ := io.ReadAll(response.Body)
			t.Fatalf("Expected 201, got %d. Body: %s", response.StatusCode, string(body))
		}

		// Проверяем ответ
		body, err := io.ReadAll(response.Body)
		if err != nil {
			t.Fatal(err)
		}

		var createdOrder order.Order
		err = json.Unmarshal(body, &createdOrder)
		if err != nil {
			t.Fatal(err)
		}

		if createdOrder.UserID != userID {
			t.Fatalf("Expected UserID %d, got %d", userID, createdOrder.UserID)
		}

		if len(createdOrder.OrderItems) != 2 {
			t.Fatalf("Expected 2 order items, got %d", len(createdOrder.OrderItems))
		}
	})

	// Тест 2: Получение всех заказов пользователя
	t.Run("GetUserOrders", func(t *testing.T) {
		req, _ := http.NewRequest("GET", testServer.URL+"/my-orders", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		client := &http.Client{}
		response, err := client.Do(req)

		if err != nil {
			t.Fatal(err)
		}

		if response.StatusCode != 200 {
			t.Fatalf("Expected 200, got %d", response.StatusCode)
		}

		body, err := io.ReadAll(response.Body)
		if err != nil {
			t.Fatal(err)
		}

		var orders []order.Order
		err = json.Unmarshal(body, &orders)
		if err != nil {
			t.Fatal(err)
		}

		if len(orders) == 0 {
			t.Fatal("Expected at least 1 order, got 0")
		}

		// Проверяем, что у пользователя есть заказ с правильным ID
		foundOrder := false
		for _, ord := range orders {
			if ord.UserID == userID {
				foundOrder = true

				// Проверяем, что OrderItems загружены
				if len(ord.OrderItems) == 0 {
					t.Fatal("OrderItems should be loaded")
				}

				// Проверяем структуру OrderItem (без лишних полей)
				for _, item := range ord.OrderItems {
					if item.ProductID == 0 {
						t.Fatal("ProductID should not be 0")
					}
					if item.Quantity <= 0 {
						t.Fatal("Quantity should be positive")
					}
				}
				break
			}
		}

		if !foundOrder {
			t.Fatalf("Order with UserID %d not found", userID)
		}
	})

	// Тест 3: Получение конкретного заказа
	t.Run("GetSpecificOrder", func(t *testing.T) {
		// Сначала получаем ID заказа
		req, _ := http.NewRequest("GET", testServer.URL+"/my-orders", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		client := &http.Client{}
		response, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}

		body, err := io.ReadAll(response.Body)
		if err != nil {
			t.Fatal(err)
		}

		var orders []order.Order
		err = json.Unmarshal(body, &orders)
		if err != nil {
			t.Fatal(err)
		}

		if len(orders) == 0 {
			t.Fatal("No orders found")
		}

		orderID := orders[0].ID

		// Теперь получаем конкретный заказ
		req, _ = http.NewRequest("GET", testServer.URL+"/order/"+strconv.Itoa(int(orderID)), nil)
		req.Header.Set("Authorization", "Bearer "+token)

		response, err = client.Do(req)
		if err != nil {
			t.Fatal(err)
		}

		if response.StatusCode != 200 {
			t.Fatalf("Expected 200, got %d", response.StatusCode)
		}

		body, err = io.ReadAll(response.Body)
		if err != nil {
			t.Fatal(err)
		}

		var specificOrder order.Order
		err = json.Unmarshal(body, &specificOrder)
		if err != nil {
			t.Fatal(err)
		}

		if specificOrder.ID != orderID {
			t.Fatalf("Expected order ID %d, got %d", orderID, specificOrder.ID)
		}

		if specificOrder.UserID != userID {
			t.Fatalf("Expected UserID %d, got %d", userID, specificOrder.UserID)
		}
	})

	// Очистка тестовых данных
	removeTestData(db)
}

func TestOrderUnauthorizedAccess(t *testing.T) {
	// Тест доступа к защищенным роутам без авторизации
	testServer := httptest.NewServer(App())
	defer testServer.Close()

	t.Run("GetOrdersWithoutAuth", func(t *testing.T) {
		req, _ := http.NewRequest("GET", testServer.URL+"/my-orders", nil)
		// Не устанавливаем Authorization header

		client := &http.Client{}
		response, err := client.Do(req)

		if err != nil {
			t.Fatal(err)
		}

		// Ожидаем 401 Unauthorized
		if response.StatusCode != 401 {
			t.Fatalf("Expected 401, got %d", response.StatusCode)
		}
	})

	t.Run("CreateOrderWithoutAuth", func(t *testing.T) {
		orderData := order.OrderCreateRequest{
			OrderItems: []order.OrderItem{
				{
					ProductID: 1,
					Quantity:  1,
				},
			},
		}

		data, _ := json.Marshal(orderData)

		req, _ := http.NewRequest("POST", testServer.URL+"/order", bytes.NewReader(data))
		req.Header.Set("Content-Type", "application/json")
		// Не устанавливаем Authorization header

		client := &http.Client{}
		response, err := client.Do(req)

		if err != nil {
			t.Fatal(err)
		}

		// Ожидаем 401 Unauthorized
		if response.StatusCode != 401 {
			t.Fatalf("Expected 401, got %d", response.StatusCode)
		}
	})
}
