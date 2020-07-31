package main

import (
	"crypto"
	"encoding/hex"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Player struct {
	board *McBoard
	id    int
}

type Lobby struct {
	board McBoard
}

var userMap = make(map[string]string)
var lobbyMap = make(map[string]*Player)
var userMapMux sync.Mutex
var lobbyMapMux sync.Mutex

func generateRoomId(username string) string {
	md5 := crypto.MD5.New()
	md5.Write([]byte(username))
	return hex.EncodeToString(md5.Sum(nil)) + fmt.Sprintf("%x", time.Now().Unix())
}

func hostGame(username string) string {
	roomId := generateRoomId(username)
	userMapMux.Lock()
	userMap[username] = roomId
	userMapMux.Unlock()
	player := new(Player)
	lobbyMapMux.Lock()
	lobbyMap[username] = player
	lobbyMapMux.Unlock()
	player.board = new(McBoard)
	*player.board = newBoard()
	player.id = 1
	return roomId
}

func joinGame(username, roomId string) error {
	return nil
}

func createLobby(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = w.Write([]byte("Invalid method"))
		return
	}
	userCookie, _ := r.Cookie("uid")
	roomId := hostGame(userCookie.Value)
	w.Write([]byte(roomId))
	return
}

func getCurrentRoom() {

}
