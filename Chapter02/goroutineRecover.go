package main

import (
	"fmt"
	"sync"
)

func simpleFunc(index int, wg *sync.WaitGroup) {
	// functions with defer keyword are executed at the end of the function
	// regardless of whether the function was executed successfully or not.
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from", r)
		}
	}()

	// We have changed the order of when wg.Done is called because
	// we should call upon wg.Done even if the following line fails.
	// Whether a defer function exists or not is dependant on whether it is registered
	// before or after the failing line of code.
	defer wg.Done()
	// This line should fail with Divide By Zero when index = 10
	fmt.Println("Attempting x/(x-10) where x = ", index, " answer is : ", index/(index-10))
}

func main() {
	var wg sync.WaitGroup
	wg.Add(40)
	for i := 0; i < 40; i += 1 {
		go func(j int) {
			simpleFunc(j, &wg)
		}(i)
	}

	wg.Wait()
}
