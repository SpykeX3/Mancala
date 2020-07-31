package main

import (
	"fmt"
	"time"
)

func log(message string) {
	println(time.Now().String(), ":", message)
	fmt.Println(message)
}
