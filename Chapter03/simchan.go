package main

import "fmt"

// helloChan waits on a channel until it gets some data and then prints the value.
func helloChan(ch <-chan string) {
	val := <-ch
	fmt.Println("Hello, ", val)
}

func main() {
	// Creating a channel
	ch := make(chan string)

	// A Goroutine that receives data from a channel
	go helloChan(ch)

	// Sending data to a channel.
	ch <- "Bob"
}
