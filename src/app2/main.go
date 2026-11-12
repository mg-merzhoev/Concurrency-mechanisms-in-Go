package main

import (
	"fmt"
	"os"
	"strconv"
)

func generate(k, n int) <-chan int {
	ch := make(chan int)

	go func() {
		defer close(ch)
		for i := k; i <= n; i++ {
			ch <- i
		}
	}()

	return ch
}

func square(in <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)
		for val := range in {

			out <- val * val
		}
	}()

	return out
}

func main() {

	if len(os.Args) < 3 {
		fmt.Println("Ошибка: необходимо передать два аргумента K и N.")
		fmt.Println("Пример запуска: go run main.go 2 5")
		return
	}

	k, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Printf("Ошибка парсинга K: %v\n", err)
		return
	}

	n, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Printf("Ошибка парсинга N: %v\n", err)
		return
	}

	// Принцип работы: конвейер(pipeline)
	chan1 := generate(k, n)
	chan2 := square(chan1)

	for res := range chan2 {
		fmt.Println(res)
	}
}
