package main

import "fmt"

const Header = "AI Training Task Processor"

func main() {
	defer fmt.Printf("%s end\n", Header)
	printHeader(Header)
	err := runScheduler(15)
	if err != nil {
		fmt.Printf("Task running error: %v\n", err)
	}
}
