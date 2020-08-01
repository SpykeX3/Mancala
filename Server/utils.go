package main

import (
	"encoding/json"
	"fmt"
	"time"
)

type errorMessage struct {
	Error string `json:"error"`
}

func log(message string) {
	println(time.Now().String(), ":", message)
	fmt.Println(message)
}

func wrapErrorJSON(err error) []byte {
	res, _ := json.Marshal(errorMessage{Error: err.Error()})
	return res
}
