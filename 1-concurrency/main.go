package main

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"time"
)

type Handler interface {
	Handle(msg string) string
}
type SimpleHandler struct{}

func (h SimpleHandler) Handle(msg string) string {
	return "Handled: " + msg
}

type HandlerWithLogging struct {
	next Handler
}

func (handlerWithLogging HandlerWithLogging) Handle(msg string) string {
	fmt.Println("Запрос: ", msg)
	result := handlerWithLogging.next.Handle(msg)
	fmt.Println("Ответ: ", result)
	return result
}

func main() {
	//---------Context with Timeout---------------
	// handler := SimpleHandler{}
	// logHandler := HandlerWithLogging{next: handler}
	// fmt.Println(logHandler.Handle("Test"))
	// ctx := context.Background()
	// contextWithTimeout, cancel := context.WithTimeout(ctx, 3*time.Second)
	// defer cancel()

	// done := make(chan struct{})

	// go func() {
	// 	time.Sleep(5 * time.Second)
	// 	close(done)
	// }()

	// select {
	// case <-done:
	// 	fmt.Println("Task is done")
	// case <-contextWithTimeout.Done():
	// 	fmt.Println("Timeout")
	// }
	//---------Context with Value-------------
	// type key int
	// const EmailKey key = 0
	// ctx := context.Background()
	// contextWithValue := context.WithValue(ctx, EmailKey, "23@gmail.com")

	// if email, ok := contextWithValue.Value(EmailKey).(string); ok {
	// 	fmt.Println(email)
	// } else {
	// 	fmt.Println("No value")
	// }
	//---------Context with Cancel-------------
	ctx, cancel := context.WithCancel(context.Background())
	go tickerOperation(ctx)
	time.Sleep(5 * time.Second)
	cancel()
	time.Sleep(5 * time.Second)

}
func tickerOperation(ctx context.Context) {
	ticker := time.NewTicker(200 * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			fmt.Println("Tick")

		case <-ctx.Done():
			fmt.Println("Canceled")
			return
		}
	}
}
func main2() {
	n := 10 // Количество генерируемых псевдослучайных чисел
	randomNumbersCh := make(chan int, n)
	squaredNumbers := make(chan int, n)

	// Генератор псведвослучайных чисел
	go func() {
		defer close(randomNumbersCh)
		for range n {
			randomNumbersCh <- rand.Intn(100)
		}
	}()

	// Обработчик возведения в квадрат
	go func() {
		defer close(squaredNumbers)
		for num := range randomNumbersCh {
			squared := int(math.Pow(float64(num), 2))
			squaredNumbers <- squared
		}
	}()

	// Вывод в строку
	for value := range squaredNumbers {
		fmt.Printf("%d ", value)
	}
}
