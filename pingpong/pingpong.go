package main

import (
	"fmt"
	"time"
)

func pong(ch chan int) {
	var i int = 0
	for i < 10 {
		i = <-ch
		fmt.Println("Pong", i)
		i++
		ch <- i
	}

}

func ping(ch chan int) {
	var i int = 0
	for i < 10 {
		fmt.Println("Ping", i)
		i++
		ch <- i
		i = <-ch
	}

}
func main() {
	ch := make(chan int)

	go ping(ch)
	go pong(ch)

	// nota, si el main termina, las goroutines acaban
	time.Sleep(1 * time.Second)
}
