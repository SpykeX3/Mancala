package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func AuthMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	notAuth := []string{"/api/user/new", "/api/user/login", "/"}
	notAuthPrefix := []string{"/js/"}
	requestPath := r.URL.Path
	println("Path is '", requestPath, "'")
	for _, value := range notAuth {
		if value == requestPath {
			println("Has special path", value)
			next(w, r)
			return
		}
	}
	for _, value := range notAuthPrefix {
		if strings.HasPrefix(requestPath, value) {
			println("Has special path prefix", value)
			next(w, r)
			return
		}
	}
	userCookie, err := r.Cookie("uid")
	if err != nil {
		log(err.Error())
		_, _ = w.Write(wrapErrorJSON(errors.New("no username was provided")))
		return
	}
	signCookie, err := r.Cookie("sign")
	if err != nil {
		log(err.Error())
		_, _ = w.Write(wrapErrorJSON(errors.New("no signature was provided")))
		return
	}
	sign, err := url.QueryUnescape(signCookie.Value)
	if err != nil {
		log(err.Error())
		_, _ = w.Write(wrapErrorJSON(errors.New("no signature was provided")))
		return
	}
	passed, err := checkAuthentication(userCookie.Value, sign)
	if !passed {
		_, _ = w.Write(wrapErrorJSON(err))
		return
	}
	fmt.Println("User: ", userCookie.Value)
	next(w, r)
}

func setUserCookie(w http.ResponseWriter, username string) {
	http.SetCookie(w, &http.Cookie{
		Name:  "uid",
		Path:  "/",
		Value: url.QueryEscape(username),
	})

	http.SetCookie(w, &http.Cookie{
		Name:  "sign",
		Path:  "/",
		Value: url.QueryEscape(sign(username)),
	})
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = w.Write(wrapErrorJSON(errors.New("invalid method")))
		return
	}
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	if username == "" || password == "" {
		_, _ = w.Write(wrapErrorJSON(errors.New("some fields were empty")))
		return
	}
	if userExists(username) {
		_, _ = w.Write(wrapErrorJSON(errors.New("user already exists")))
		return
	}
	addUser(username, password)
	setUserCookie(w, username)
	println("Added user:", username, password)
}

func loginUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = w.Write(wrapErrorJSON(errors.New("invalid method")))
		return
	}
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	if !checkPassword(username, password) {
		_, _ = w.Write(wrapErrorJSON(errors.New("invalid credentials")))
		return
	}
	setUserCookie(w, username)
	println("Logged in as", username, password)
}
