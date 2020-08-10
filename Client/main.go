package main

import (
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"log"
	"net/url"
	"runtime"
	"strings"
	"syscall"
	"time"
)

func panicCheck(err error) {
	if err != nil {
		log.Panicln(err)
	}
}
func readString() string {
	userStr := ""
	for userStr == "" {
		_, err := fmt.Scan(&userStr)
		if err != nil {
			log.Panicln("Failed to read url")
		}
	}
	return userStr
}

func readIntInRange(first, count int) int {
	cmd := first - 1
	if count < 1 {
		log.Panicln("Invalid count")
	}
	for cmd < first || cmd > first+count-1 {
		n, err := fmt.Scan(&cmd)
		if err != nil {
			log.Println("Failed to read command")
		}
		if n == 0 {
			continue
		}
	}
	return cmd
}

func readURL() string {
	urlStr := ""
	for urlStr == "" {
		urlStr = readString()
		_, err := url.Parse(urlStr)
		if err != nil {
			log.Println("Invalid URL was provided")
			urlStr = ""
			continue
		}
	}
	if !strings.HasSuffix(urlStr, "/") {
		urlStr += "/"
	}
	return urlStr
}

func readCredentials() (string, string) {
	username, password := "", ""
	fmt.Println("Username:")
	username = readString()
	fmt.Println("Password:")
	if runtime.GOOS == "windows" {
		password = readString()
	} else {
		passwordBytes, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			log.Panicln("Failed to read password")
		}
		password = string(passwordBytes)
	}
	return username, password
}

func main() {
	var username string
	fmt.Println("Please, type in the server URL")
	fmt.Println("example:\thttp://localhost:1234/")
	urlStr := readURL()
	user := newMancalaClient(urlStr)
	var cmd int
	fmt.Println("Welcome to the Game of Mancala!")
	for !user.IsLoggedIn() {
		fmt.Println("Type 1 to sing in or 2 to sign up")
		cmd = readIntInRange(1, 2)
		var password string
		username, password = readCredentials()
		switch cmd {
		case 1:
			if err := user.SignIn(username, password); err != nil {
				fmt.Println(err.Error())
			}
		case 2:
			if err := user.SignUp(username, password); err != nil {
				fmt.Println(err.Error())
			}
		}
	}
	fmt.Println("Successfully logged in")
	for true {
		// Join a game
		for !user.IsInLobby() {
			fmt.Println("Type 1 to host a game, 2 to join")
			cmd = readIntInRange(1, 2)
			switch cmd {
			case 1:
				room, err := user.CreateLobby()
				if err != nil {
					fmt.Println(err.Error())
				}
				fmt.Println("Your room id is:")
				fmt.Println(room)
			case 2:
				fmt.Println("Type room id to join")
				room := readString()
				if err := user.JoinLobby(room); err != nil {
					fmt.Println(err.Error())
				}
			}
		}
		// Game loop
		fmt.Println("Hint:\tOn your turn type 1-6 ")
		lastUpdatedState := 0
		for true {
			// Get board state
			board, err := user.GetGameState() //TODO show if game has changed
			if err != nil {
				fmt.Println(err.Error())
				time.Sleep(time.Second)
				continue
			}
			// Exit loop if game is over
			if board.Result.GameOver {
				fmt.Print("The game is over:\t")
				if board.Result.IsDraw {
					fmt.Println("Draw")
				} else {
					fmt.Println(board.Result.Winner, "wins")
				}
				fmt.Println(board.Printable(username))
				user.LeaveLobby()
				break
			}
			//Show state if it was updated
			if board.LastChanged != lastUpdatedState {
				fmt.Println(board.Printable(username))
				lastUpdatedState = board.LastChanged
			}
			// If user should make a turn, go into the input loop
			if board.Players[board.NextPlayer-1] == username {
				fmt.Println("Your turn:")
				for true {
					turn := readIntInRange(1, 6)
					turn--
					board, err = user.MakeTurn(turn)
					if err == nil {
						break
					}
					fmt.Println(err.Error())
					time.Sleep(time.Second)
					fmt.Println("Try again")
				}
			}
			time.Sleep(time.Second)
		}
	}
}
