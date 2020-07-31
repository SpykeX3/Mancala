package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"math/rand"
	"sync"
)

var users = make(map[string]string)
var mutx sync.Mutex

func generateKey() string {
	result := ""
	for i := 0; i < 4; i++ {
		result += fmt.Sprintf("%x", rand.Int63())
	}
	return result
}

func checkAuthentication(username, signature string) (bool, error) {
	if !userExists(username) {
		return false, errors.New("corrupted cookie")
	}
	hmc := hmac.New(sha256.New, CookieKey)
	hmc.Write([]byte(username))
	valid := hmc.Sum(nil)

	decodedSignature, err := base64.URLEncoding.DecodeString(signature)
	if err != nil {
		return false, err
	}
	if !hmac.Equal(valid, decodedSignature) {
		println(valid)
		println(decodedSignature)
		return false, errors.New("corrupted cookie")
	}
	return true, nil
}

func sign(message string) string {
	hmc := hmac.New(sha256.New, CookieKey)
	hmc.Write([]byte(message))
	sum := hmc.Sum(nil)
	return base64.URLEncoding.EncodeToString(sum)
}

func checkPassword(username, password string) bool {
	mutx.Lock()
	//TODO real DB
	realPass, found := users[username]
	mutx.Unlock()
	if !found {
		return false
	}
	return realPass == password
}

func userExists(username string) bool {
	mutx.Lock()
	_, exists := users[username]
	mutx.Unlock()
	return exists
}
func addUser(username, password string) {
	mutx.Lock()
	users[username] = password
	defer mutx.Unlock()
}
