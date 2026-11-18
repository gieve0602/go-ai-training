package main

import (
	"errors"
	"fmt"
)

func normal() int {
	n := 0
	defer func() {
		n++
	}()
	for range 5 {
		n++
	}
	return n
}

func named() (n int) {
	defer func() {
		n++
	}()
	for range 5 {
		n++
	}
	return
}

type (
	operator      func(x, y int) int
	calculateFunc func(x, y int) (int, error)
)

func genCalculator(op operator) calculateFunc {
	return func(x, y int) (int, error) {
		if op == nil {
			return 0, errors.New("invalid operator")
		}
		return op(x, y), nil
	}
}

func modifyArray(arr [3][]string) {
	arr[1][1] = "h1"
}

func main() {
	n1 := normal()
	n2 := named()

	fmt.Printf("normal: %d\n", n1)
	fmt.Printf("named: %d\n", n2)

	op := func(x, y int) int {
		return x + y
	}
	add := genCalculator(op)
	result, err := add(5, 10)
	if err != nil {
		fmt.Printf("operator add error: %v\n", err)
	}
	fmt.Printf("add result: %d\n", result)

	complexArray1 := [3][]string{
		{"d", "e", "f"},
		{"g", "h", "i"},
		{"j", "k", "l"},
	}
	fmt.Printf("original arr: %#v\n", complexArray1)
	modifyArray(complexArray1)
	fmt.Printf("modified arr: %#v\n", complexArray1)
}
