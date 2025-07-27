package main

import (
	"fmt"
	"math"
	"math/rand"
)

func main() {
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
