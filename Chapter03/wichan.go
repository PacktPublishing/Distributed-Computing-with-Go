package main

import (
	"fmt"
	"sync"
)

func cashier(cashierID int, orderChannel <-chan int, wg *sync.WaitGroup) {
	// Process orders upto limit.
	for ordersProcessed := 0; ordersProcessed < 10; ordersProcessed++ {
		// Retrieve order from orderChannel
		orderNum := <-orderChannel

		// Cashier is ready to serve!
		fmt.Println("Cashier ", cashierID, "Processing order", orderNum, "Orders Processed", ordersProcessed)
		wg.Done()
	}
}

func main() {
	var wg sync.WaitGroup
	wg.Add(30)
	ordersChannel := make(chan int)

	for i := 0; i < 3; i++ {
		// Start the three cashiers
		func(i int) {
			go cashier(i, ordersChannel, &wg)
		}(i)
	}

	// Start adding orders to be processed.
	for i := 0; i < 30; i++ {
		ordersChannel <- i
	}
	wg.Wait()
}
