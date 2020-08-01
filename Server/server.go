package main

import (
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

func main() {
	initFromEnv()
	router := mux.NewRouter()
	router.HandleFunc("/api/user/new", createUserHandler)
	router.HandleFunc("/api/user/login", loginUserHandler)
	router.HandleFunc("/api/lobby/create", createLobbyHandler)
	router.HandleFunc("/api/lobby/join", joinLobbyHandler)
	router.HandleFunc("/api/lobby/turn", makeTurnHandler)
	n := negroni.New()
	n.UseHandler(negroni.Classic())
	n.UseFunc(AuthMiddleware)
	n.UseHandler(router)
	err := http.ListenAndServe(":"+port, n)
	if err != nil {
		panic(err.Error())
	}
}
