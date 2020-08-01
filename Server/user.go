package main

import (
	"fmt"
	"net/http"
	"time"
)

func AuthMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	notAuth := []string{"/api/user/new", "/api/user/login"} //Список эндпоинтов, для которых не требуется авторизация
	requestPath := r.URL.Path                               //текущий путь запроса
	for _, value := range notAuth {

		if value == requestPath {
			next(w, r)
			return
		}
	}
	userCookie, err := r.Cookie("uid")
	if err != nil {
		log(err.Error())
		_, _ = w.Write([]byte("No username was provided"))
		return
	}
	signCookie, err := r.Cookie("sign")
	if err != nil {
		log(err.Error())
		_, _ = w.Write([]byte("No signature was provided"))
		return
	}
	passed, err := checkAuthentication(userCookie.Value, signCookie.Value)
	if !passed {
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	fmt.Println("User: ", userCookie.Value)
	next(w, r)
}

func setUserCookie(w http.ResponseWriter, username string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "uid",
		Value:    username,
		Path:     "/",
		Expires:  time.Now().Add(time.Hour * 24),
		HttpOnly: true,
		SameSite: 1,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "sign",
		Value:    sign(username),
		Path:     "/",
		Expires:  time.Now().Add(time.Hour * 24),
		HttpOnly: true,
		SameSite: 1,
	})
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = w.Write([]byte("Invalid method"))
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
		_, _ = w.Write([]byte("Some fields were empty"))
	}
	if userExists(username) {
		_, _ = w.Write([]byte("User already exists"))
	}
	addUser(username, password)
	setUserCookie(w, username)
}

func loginUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = w.Write([]byte("Invalid method"))
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
		_, _ = w.Write([]byte("Invalid credentials"))
		return
	}
	setUserCookie(w, username)
}
