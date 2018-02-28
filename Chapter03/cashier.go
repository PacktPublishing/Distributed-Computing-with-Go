package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	// ordersProcessed & cashier are declared in main function
	// so that cashier has access to shared state variable `ordersProcessed`.
	// If we were to declare the variable inside the `cashier` function,
	// then it's value would be set to zero with every function call.
	ordersProcessed := 0
	cashier := func(orderNum int) {
		if ordersProcessed < 10 {
			// Cashier is ready to serve!
			fmt.Println("Processing order", orderNum)
			ordersProcessed++
		} else {
			// Cashier has reached the max capacity of processing orders.
			fmt.Println("I am tired! I want to take rest!", orderNum)
		}
		wg.Done()
	}

	for i := 0; i < 30; i++ {
		// Note that instead of wg.Add(60), we are instead adding 1
		// per each loop iteration. Both are valid ways to add to WaitGroup as long as we can ensure the right number of calls.
		wg.Add(1)
		go func(orderNum int) {
			// Making an order
			cashier(orderNum)
		}(i)

	}
	wg.Wait()
}
