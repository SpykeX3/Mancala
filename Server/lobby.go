package main

import (
	"Mancala/Server/Mancala"
	"crypto"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
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
	board            *Mancala.McBoard
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

	lobby := &Lobby{
		name: "Lobby of " + username,
		player1: &Player{
			connection: p1Connection,
			username:   "",
			id:         1,
		},
		player2:   nil,
		board:     Mancala.NewBoard(),
		stateChan: make(chan string),
	}
	lobby.board.Players[0] = username
	lobbyMapMux.Lock()
	lobbyMap[roomId] = lobby
	lobbyMapMux.Unlock()
	return roomId
}

func joinGame(username, roomId string) error {
	lobbyMapMux.Lock()
	defer lobbyMapMux.Unlock()
	lobby, exists := lobbyMap[roomId]
	if !exists {
		return errors.New("invalid room identifier")
	}
	if lobby.player2 != nil {
		return errors.New("lobby is full")
	}

	userMapMux.Lock()
	currentRoom, inRoom := userMap[username]
	if inRoom && currentRoom == roomId {
		userMapMux.Unlock()
		return errors.New("user is already in room")
	}
	userMap[username] = roomId
	userMapMux.Unlock()

	connectionMapMux.Lock()
	p2Connection := newGCConnection()
	connectionMap[username] = p2Connection
	connectionMapMux.Unlock()
	lobby.player2 = &Player{
		connection: p2Connection,
		username:   username,
		id:         2,
	}
	lobby.board.Players[1] = username
	go gameControllerRoutine(lobby.board, lobby.player1.connection, lobby.player2.connection, &lobby.stateChan)
	return nil
}

func gameControllerRoutine(board *Mancala.McBoard, player1, player2 *gcConnection, stateChan *chan string) {
	killer := make(chan bool, 0)
	killOnce := sync.Once{}
	if board.Result.GameOver {
		time.AfterFunc(time.Minute*60, func() {
			killer <- true
		})
	}
	for {
		select {
		case p1turn := <-player1.request:
			{
				err := board.Turn(1, p1turn)
				if err != nil {
					player1.response <- string(wrapErrorJSON(err))
				} else {
					jsOut, _ := json.Marshal(*board)
					player1.response <- string(jsOut)
				}

			}
		case p2turn := <-player2.request:
			{
				err := board.Turn(2, p2turn)
				if err != nil {
					player2.response <- string(wrapErrorJSON(err))
				} else {
					jsOut, _ := json.Marshal(*board)
					player2.response <- string(jsOut)
				}
			}
		case *stateChan <- board.String():
			{
			}
		case <-killer:
			{
				log.Println("Leaving gameController goroutine for users", board.Players)
				return
			}
		}
		if board.Result.GameOver {
			killOnce.Do(func() {
				time.AfterFunc(time.Second*10, func() {
					//log.Println("Game is finished, sending killing signal")
					killer <- true
				})
			})
		}
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
	select {
	case connection.request <- cell:
		return <-connection.response
	default:
		err := errors.New("game is not active")
		return string(wrapErrorJSON(err))
	}
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
	select {
	case resp := <-lobby.stateChan:
		return resp
	default:
		err := errors.New("game is not active")
		return string(wrapErrorJSON(err))
	}
}
