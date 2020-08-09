package main

import (
	"encoding/json"
)

type errorMessage struct {
	Error string `json:"error"`
}

func wrapErrorJSON(err error) []byte {
	res, _ := json.Marshal(errorMessage{Error: err.Error()})
	return res
}
