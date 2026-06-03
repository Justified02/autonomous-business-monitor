package main

import "fmt"

func calculate(ch chan int) {
	ch <- 42
}

func double(ch chan int) {
	ch <- 84
}

func triple(ch chan int) {
	ch <- 126
}

func main() {
	ch := make(chan int, 3)
	go calculate(ch)
	go double(ch)
	go triple(ch)

	result1 := <- ch
	result2 := <- ch
	result3 := <- ch

	fmt.Println("Got:", result1)
	fmt.Println("Got:", result2)
	fmt.Println("Got:", result3)
}