package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println(os.Getenv("NAME") + " is your uncle.")
}
