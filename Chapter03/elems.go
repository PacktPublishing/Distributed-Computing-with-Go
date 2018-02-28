package main

import "fmt"

func main() {
	// Let's create three simple functions that take an int argument
	fcn1 := func(i int) {
		fmt.Println("fcn1", i)
	}
	fcn2 := func(i int) {
		fmt.Println("fcn2", i*2)
	}
	fcn3 := func(i int) {
		fmt.Println("fcn3", i*3)
	}

	ch := make(chan func(int)) // Channel that sends & receives functions that take an int argument
	done := make(chan bool)    // A Channel whose element type is a boolean value.

	// Launch a goroutine to work with the channels ch & done.
	go func() {
		// We accept all incoming functions on Channel ch and call the functions with value 10.
		for fcn := range ch {
			fcn(10)
		}
		// Once the loop terminates, we print Exiting and send true to done Channel.
		fmt.Println("Exiting")
		done <- true
	}()

	// Sending functions to channel ch
	ch <- fcn1
	ch <- fcn2
	ch <- fcn3

	// Close the channel once we are done sending it data.
	close(ch)

	// Wait on the launched goroutine to end.
	<-done
}
