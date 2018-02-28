package common

import (
	"fmt"
	"log"
)

func Log(msg string) {
	log.Println("INFO - ", msg)
}

func Warn(msg string) {
	log.Println("---------------------------")
	log.Println(fmt.Sprintf("WARN: %s", msg))
	log.Println("---------------------------")
}
