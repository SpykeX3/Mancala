package main

import (
	"crypto"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"
)

type Player struct {
	connection *gcConnection
	username   string
	id         int
}

type gcConnection struct {
	request  chan int
	response chan string
}

type Lobby struct {
	name             string
	board            *McBoard
	player1, player2 *Player
	stateChan        chan string
}

//###
//TODO Not cool :(
var userMap = make(map[string]string)
var lobbyMap = make(map[string]*Lobby)
var connectionMap = make(map[string]*gcConnection)
var userMapMux sync.Mutex
var lobbyMapMux sync.Mutex
var connectionMapMux sync.Mutex

//####

func generateRoomId(username string) string {
	md5 := crypto.MD5.New()
	md5.Write([]byte(username))
	return hex.EncodeToString(md5.Sum(nil)) + fmt.Sprintf("%x", time.Now().Unix())
}

func newGCConnection() *gcConnection {
	return &gcConnection{
		request:  make(chan int),
		response: make(chan string),
	}
}

func hostGame(username string) string {
	roomId := generateRoomId(username)
	userMapMux.Lock()
	userMap[username] = roomId
	userMapMux.Unlock()

	connectionMapMux.Lock()
	p1Connection := newGCConnection()
	connectionMap[username] = p1Connection
	connectionMapMux.Unlock()

	lobbyMapMux.Lock()
	lobbyMap[roomId] = &Lobby{
		name: "Lobby of " + username,
		player1: &Player{
			connection: p1Connection,
			username:   "",
			id:         1,
		},
		player2:   nil,
		board:     newBoard(),
		stateChan: make(chan string),
	}
	lobbyMapMux.Unlock()
	return roomId
}

func joinGame(username, roomId string) error {
	lobbyMapMux.Lock()
	lobby, exists := lobbyMap[roomId]
	if !exists {
		lobbyMapMux.Unlock()
		return errors.New("invalid room identifier")
	}
	if lobby.player2 != nil {
		lobbyMapMux.Unlock()
		return errors.New("lobby is full")
	}
	connectionMapMux.Lock()
	p2Connection := newGCConnection()
	connectionMap[username] = p2Connection
	connectionMapMux.Unlock()
	lobby.player2 = &Player{
		connection: p2Connection,
		username:   username,
		id:         2,
	}
	lobbyMapMux.Unlock()
	userMapMux.Lock()
	userMap[username] = roomId
	userMapMux.Unlock()
	go gameControllerRoutine(lobby.board, lobby.player1.connection, lobby.player2.connection, &lobby.stateChan)
	return nil
}

func gameControllerRoutine(board *McBoard, player1, player2 *gcConnection, stateChan *chan string) {
	for {
		select {
		case p1turn := <-player1.request:
			{
				err := board.turn(1, p1turn)
				if err != nil {
					player1.response <- string(wrapErrorJSON(err))
				} else {
					jsOut, _ := json.Marshal(*board)
					player1.response <- string(jsOut)
				}

			}
		case p2turn := <-player2.request:
			{
				err := board.turn(2, p2turn)
				if err != nil {
					player2.response <- string(wrapErrorJSON(err))
				} else {
					jsOut, _ := json.Marshal(*board)
					player2.response <- string(jsOut)
				}
			}
		case *stateChan <- board.string():
			{
			}
		}
		//TODO finite game \o/
	}
}

func makeTurn(username string, cell int) string {
	connectionMapMux.Lock()
	defer connectionMapMux.Unlock()
	connection, exists := connectionMap[username]
	if !exists {
		err := errors.New("user is not in any game")
		return string(wrapErrorJSON(err))
	}
	connection.request <- cell
	return <-connection.response
}

func getGameState(username string) string {
	userMapMux.Lock()
	room, exists := userMap[username]
	userMapMux.Unlock()
	if !exists {
		err := errors.New("user is not in any game")
		return string(wrapErrorJSON(err))
	}
	lobbyMapMux.Lock()
	lobby, exists := lobbyMap[room]
	lobbyMapMux.Unlock()
	if !exists {
		err := errors.New("user is not in any game")
		return string(wrapErrorJSON(err))
	}
	return <-lobby.stateChan
}
