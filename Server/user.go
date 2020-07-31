package main

import (
	"net/http"
	"time"
)

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

func createUser(w http.ResponseWriter, r *http.Request) {
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

func loginUser(w http.ResponseWriter, r *http.Request) {
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
