package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/urfave/negroni"
	"net/http"
	"os"
)

var CookieKey []byte
var port string

func initFromEnv() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	println(wd)
	err = godotenv.Load(".env")
	if err != nil {
		println(err.Error())
	}
	key := os.Getenv("MANCALA_KEY")
	port = os.Getenv("MANCALA_PORT")
	if key == "" {
		key = generateKey()
		keyMap := make(map[string]string)
		keyMap["MANCALA_KEY"] = key
		err = godotenv.Write(keyMap, ".env")
		if err != nil {
			println(err.Error())
		}
	}
	if port == "" {
		port = "8080"
	}
	CookieKey = []byte(key)
	println("KEY:", key)
	println("PORT:", port)
}

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

func main() {
	initFromEnv()
	router := mux.NewRouter()
	router.HandleFunc("/api/user/new", createUser)
	router.HandleFunc("/api/user/login", loginUser)
	router.HandleFunc("/api/lobby/create", createLobby)
	n := negroni.New()
	n.UseHandler(negroni.Classic())
	n.UseFunc(AuthMiddleware)
	n.UseHandler(router)
	err := http.ListenAndServe(":"+port, n)
	if err != nil {
		panic(err.Error())
	}
}
