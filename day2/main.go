package main

import (
	"fmt"
	"unicode/utf8"
)

const (
	SEP   string = "================================"
	TITLE string = "AI Training Configuration"
)

const (
	Pending int = iota
	Running
	Success
	Failed
)

var (
	statusUnknown string = "Unknown"
	statusNames          = map[int]string{
		Pending: "Pending",
		Running: "Running",
		Success: "Success",
		Failed:  "Failed",
	}
)

func main() {
	application := "AI Training"
	version := "1.0.0"
	port := 8080
	maxConnections := 100
	debug := true
	timeout := 30
	status := 1
	fmt.Println(SEP)
	fmt.Println(TITLE)
	fmt.Println(SEP)
	fmt.Printf("Application: %v\n", application)
	fmt.Printf("Version: %v\n", version)
	fmt.Printf("Port: %v\n", port)
	fmt.Printf("MaxConnections: %v\n", maxConnections)
	fmt.Printf("Debug: %v\n", debug)
	fmt.Printf("Timeout: %v\n", timeout)
	if statusName, ok := statusNames[status]; ok {
		fmt.Printf("Status: %v\n", statusName)
	} else {
		fmt.Printf("Status: %v\n", statusUnknown)
	}
	fmt.Println(SEP)
	fmt.Println("Type Information")
	fmt.Printf("Port: %T\n", port)
	fmt.Printf("Version: %T\n", version)
	fmt.Printf("Debug: %T\n", debug)
	fmt.Println(SEP)
	fmt.Println("Extension excrise")
	var tmp uint8 = 255
	fmt.Printf("uint 255: original=%v, plus one=%v\n", tmp, tmp+1)
	tmpStr := "你好"
	fmt.Printf("chinese string lenth: %v, runeCount: %v\n", len(tmpStr), utf8.RuneCountInString(tmpStr))
	var zeroInt int
	var zeroStr string
	var zeroBool bool
	var zeroMap map[string]int
	fmt.Printf("zero value: int=%v, string=%v, bool=%v, map=%v, mapIsNil=%v\n", zeroInt, zeroStr, zeroBool, zeroMap, zeroMap == nil)
	var a1 uint8 = 200
	a2 := float64(a1)
	fmt.Printf("transform: a1=%v, a2=%v\n", a1, a2)
	fmt.Println(SEP)
}
