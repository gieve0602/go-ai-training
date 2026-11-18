package main

import (
	"fmt"
)

type User struct {
	Name string
	Age  int
	Male bool
}

func main() {
	//user := &User{
	//	Name: "test1",
	//	Age:  19,
	//	Male: true,
	//}

	var user1 *User
	if user1 != nil && user1.Age >= 18 {
		fmt.Println("Adult")
	}

	r1 := 5 + 16*2
	r2 := (5 + 16) * 2
	r3 := 5 + (16 * 2)
	fmt.Printf("exp1: 5 + 16 * 2 = %v\n", r1)
	fmt.Printf("exp2: (5 + 16) * 2 = %v\n", r2)
	fmt.Printf("exp3: 5 + (16 * 2) = %v\n", r3)

	read := 1
	write := 2
	execute := 4

	readWrite := read | write
	rwx := read | write | execute
	fmt.Println(rwx)

	if readWrite&read == read {
		fmt.Println("has read permission")
	}
	if readWrite&write == write {
		fmt.Println("has write permission")
	}
	if rwx&execute == execute {
		fmt.Println("has execute permission")
	}
	readWrite = readWrite &^ write
	if readWrite&write != write {
		fmt.Println("has removed write")
	}

	fmt.Println("")
	for i := 1; i <= 20; i++ {
		if i%3 == 0 {
			continue
		}
		if i == 17 {
			break
		}
		fmt.Printf("num: %v\n", i)
	}

	tmpStr := "Go语言工程"
	for idx, c := range tmpStr {
		fmt.Printf("rune idx=%d, char=%c\n", idx, c)
	}
}
