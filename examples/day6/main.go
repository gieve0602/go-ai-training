package main

import "fmt"

func modify(s *[]int) {
	(*s)[0] = 100
	*s = append(*s, 200)
}

func change(p *int) {
	p = new(int)
	*p = 100
}

type Runner interface {
	Run()
}

type Task struct{}

func (t *Task) Run() {}

func main() {
	var c []int
	var d map[string]int
	fmt.Println(c == nil)
	fmt.Println(d == nil)
	nums := []int{1, 2, 3}
	modify(&nums)
	fmt.Println(nums)

	x := 10
	change(&x)
	fmt.Println(x)

	var a Runner = Task{}
	var b Runner = &Task{}
	_ = a
	_ = b
}
