package main

import (
	"fmt"
	"sync"
)

func simpleFunc(index int, wg *sync.WaitGroup) {
	// This line should fail with Divide By Zero when index = 10
	fmt.Println("Attempting x/(x-10) where x = ", index, " answer is : ", index/(index-10))
	wg.Done()
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
