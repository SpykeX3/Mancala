package main

import (
	"io/ioutil"
	"net/http"
	"strconv"
)

func createLobbyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = w.Write([]byte("Invalid method"))
		return
	}
	userCookie, _ := r.Cookie("uid")
	roomId := hostGame(userCookie.Value)
	_, _ = w.Write([]byte(roomId))
	return
}

func joinLobbyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = w.Write([]byte("Invalid method"))
		return
	}
	err := r.ParseForm()
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	roomId := r.Form.Get("room")
	userCookie, _ := r.Cookie("uid")
	err = joinGame(userCookie.Value, roomId)
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	_, _ = w.Write([]byte("Connected to " + roomId))
	return
}

func makeTurnHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = w.Write([]byte("Invalid method"))
		return
	}
	err := r.ParseForm()
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	userCookie, _ := r.Cookie("uid")

	cellStr := r.Form.Get("cell")
	cell, err := strconv.Atoi(cellStr)
	if err != nil {
		_, _ = w.Write(wrapErrorJSON(err))
		return
	}
	_, _ = w.Write([]byte(makeTurn(userCookie.Value, cell)))
	return
}

func gameStateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = w.Write([]byte("Invalid method"))
		return
	}
	userCookie, _ := r.Cookie("uid")
	_, _ = w.Write([]byte(getGameState(userCookie.Value)))
	return
}

func mainPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = w.Write([]byte("Invalid method"))
		return
	}
	contents, err := ioutil.ReadFile("index.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log("Main page was inaccessible")
		return
	}
	_, _ = w.Write(contents)
	return
}
