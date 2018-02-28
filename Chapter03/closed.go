package main

import "fmt"

type msg struct {
	ID    int
	value string
}

func handleIntChan(intChan <-chan int, done chan<- int) {
	// Even though there are only 4 elements being sent via channel, we retrieve 6 values.
	for i := 0; i < 6; i++ {
		fmt.Println(<-intChan)
	}
	done <- 0
}

func handleMsgChan(msgChan <-chan msg, done chan<- int) {
	// We retrieve 6 values of element type struct `msg`.
	// Given that there are only 4 values in the buffered channel,
	// the rest should be zero value of struct `msg`.
	for i := 0; i < 6; i++ {
		fmt.Println(fmt.Sprintf("%#v", <-msgChan))
	}
	done <- 0
}

func main() {
	intChan := make(chan int)
	done := make(chan int)

	go func() {
		intChan <- 9
		intChan <- 2
		intChan <- 3
		intChan <- 7
		close(intChan)
	}()
	go handleIntChan(intChan, done)

	msgChan := make(chan msg, 5)
	go func() {
		for i := 1; i < 5; i++ {
			msgChan <- msg{
				ID:    i,
				value: fmt.Sprintf("VALUE-%v", i),
			}
		}
		close(msgChan)
	}()
	go handleMsgChan(msgChan, done)

	// We wait on the two channel handler goroutines to complete.
	<-done
	<-done

	// Since intChan is closed, this will cause a panic to occur.
	intChan <- 100
}
