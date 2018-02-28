package main

import (
	"fmt"
	"sync"
)

func recv(ch <-chan int, wg *sync.WaitGroup) {
	fmt.Println("Receiving", <-ch)
	wg.Done()
}

func send(ch chan<- int, wg *sync.WaitGroup) {
	fmt.Println("Sending...")
	ch <- 100
	fmt.Println("Sent")
	wg.Done()
}

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	ch := make(chan int)
	go recv(ch, &wg)
	go send(ch, &wg)

	wg.Wait()
}
