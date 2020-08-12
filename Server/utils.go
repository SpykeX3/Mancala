package main

import (
	"encoding/json"
	"unicode"
)

type errorMessage struct {
	Error string `json:"error"`
}

func capitalizeFirst(str string) string {
	runes := []rune(str)
	unicode.ToUpper(runes[0])
	return string(runes)
}
func wrapErrorJSON(err error) []byte {
	res, _ := json.Marshal(errorMessage{Error: capitalizeFirst(err.Error())})
	return res
}

func newErrorJSON(message string) []byte {
	res, _ := json.Marshal(errorMessage{Error: capitalizeFirst(message)})
	return res
}
