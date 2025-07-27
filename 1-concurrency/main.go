package main

import (
	"fmt"
	"math"
	"math/rand"
)

func main() {
	n := 10
	randomNumbersCh := make(chan int, n)
	squaredNumbers := make(chan int, n)

	go func() {
		defer close(randomNumbersCh)
		for range n {
			randomNumbersCh <- rand.Intn(100)
		}
	}()

	go func() {
		defer close(squaredNumbers)
		for num := range randomNumbersCh {
			squared := int(math.Pow(float64(num), 2))
			squaredNumbers <- squared
		}
	}()

	for value := range squaredNumbers {
		fmt.Println(value)
	}
}
