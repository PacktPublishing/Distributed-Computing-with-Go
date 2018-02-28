package main

import "fmt"

func main() {
	channels := [5](chan int){
		make(chan int),
		make(chan int),
		make(chan int),
		make(chan int),
		make(chan int),
	}

	go func() {
		// Starting to wait on channels
		for _, chX := range channels {
			fmt.Println("Receiving from", <- chX)
		}
	}()

	for i := 1; i < 6; i++ {
		fmt.Println("Sending on channel:", i)
		channels[i] <- 1
	}
}
