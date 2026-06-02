package main

import "fmt"

func calculate(ch chan int) {
	ch <- 42
}

func double(ch chan int) {
	ch <- 84
}

func main() {
	ch := make(chan int, 1)
	che := make(chan int, 1)
	go calculate(ch)
	go double(che)
	result := <-che
	fmt.Println("Got:", result)
}