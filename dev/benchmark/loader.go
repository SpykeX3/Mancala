package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"sync"
	"time"
)

type errorMessage struct {
	Error string `json:"error"`
}

type GameResult struct {
	GameOver bool   `json:"game_over"`
	IsDraw   bool   `json:"is_draw"`
	Winner   string `json:"winner"`
}

type McCell struct {
	Score int `json:"score"`
}

type McBoard struct {
	P1score     int        `json:"player1_score"`
	P2score     int        `json:"player2_score"`
	P1cells     []McCell   `json:"player1_cells"`
	P2cells     []McCell   `json:"player2_cells"`
	P1mc        McCell     `json:"player1_mancala"`
	P2mc        McCell     `json:"player2_mancala"`
	NextPlayer  int        `json:"next_player"`
	Players     [2]string  `json:"players"`
	Result      GameResult `json:"result"`
	LastChanged int        `json:"last_changed"`
}

func client(url string, id int, interval time.Duration, lobbies chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	username := fmt.Sprintf("testingAccount%d", id)
	password := fmt.Sprintf("testingPassword%d", id)
	var room string
	jar, _ := cookiejar.New(nil)
	client := http.Client{
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           jar,
		Timeout:       0,
	}
	_, _ = client.Post(url+"/api/user/new", "application/x-www-form-urlencoded", strings.NewReader(fmt.Sprintf("username=%s&password=%s", username, password)))
	_, _ = client.Post(url+"/api/user/login", "application/x-www-form-urlencoded", strings.NewReader(fmt.Sprintf("username=%s&password=%s", username, password)))
	if id&1 == 0 {
		resp, _ := client.Post(url+"/api/lobby/create", "application/x-www-form-urlencoded", nil)
		respB, _ := ioutil.ReadAll(resp.Body)
		room = string(respB)
		log.Printf("[%d] Hosting %s\n", id, room)
		lobbies <- room
	} else {
		room = <-lobbies
		_, _ = client.Post(url+"/api/lobby/join", "application/x-www-form-urlencoded", strings.NewReader(fmt.Sprintf("room=%s", room)))
		log.Printf("[%d] Joining %s\n", id, room)
	}
	cell := 3
	for true {
		time.Sleep(interval)
		//log.Printf("[%d] Next cell is %d\n", id, cell)
		resp, _ := client.Post(url+"/api/lobby/turn", "application/x-www-form-urlencoded", strings.NewReader(fmt.Sprintf("cell=%d", cell)))
		respB, _ := ioutil.ReadAll(resp.Body)
		//log.Printf("[%d] Got message %s\n", id, string(respB))
		var board McBoard
		var respError errorMessage
		dec := json.NewDecoder(bytes.NewReader(respB))
		dec.DisallowUnknownFields()
		err := dec.Decode(&respError)
		if err == nil {
			if respError.Error == "game is over" {
				log.Println("Finished ", room)
				break
			}
		}
		err = json.Unmarshal(respB, &board)
		cell = (cell + 1) % 6
		if err != nil {
			log.Panicln(err)
		}
		if board.Result.GameOver {
			log.Println("Finished ", room)
			break
		}
	}
}
func main() {
	count := 200
	url := "http://localhost:1337"
	var wg sync.WaitGroup
	matchmaker := make(chan string, count)
	wg.Add(count)
	start := time.Now()
	for i := 0; i < count; i++ {
		go client(url, i, time.Millisecond*250, matchmaker, &wg)
	}
	wg.Wait()
	println("Time:", time.Now().Sub(start).String())
}
