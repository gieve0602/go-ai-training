package main

import "fmt"

type CalculateFunc func(int, int) int

func printHeader(header string) {
	fmt.Println("=========================")
	fmt.Println(header)
	fmt.Println("=========================")
}

func sum(numbers ...int) int {
	sum := 0
	for _, n := range numbers {
		sum += n
	}
	return sum
}

func add(x, y int) int {
	return x + y
}

func subtract(x, y int) int {
	return x - y
}

func multiply(x, y int) int {
	return x * y
}

func execute(x, y int, fn CalculateFunc) int {
	return fn(x, y)
}

func createTaskCounter() func() int {
	count := 0
	return func() int {
		count++
		return count
	}
}
