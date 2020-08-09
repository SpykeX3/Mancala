package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

type errorMessage struct {
	Error string `json:"error"`
}

type MancalaClient interface {
	SignUp(username, password string) error
	SignIn(username, password string) error
	CreateLobby() (string, error)
	JoinLobby(room string) error
	GetGameState() (*McBoard, error)
	MakeTurn(cell int) (*McBoard, error)
	IsLoggedIn() bool
	IsInLobby() bool
	LeaveLobby()
}

type MancalaClientProt struct {
	url      string
	loggedIn bool
	inLobby  bool
	client   *http.Client
}

func (mc MancalaClientProt) IsLoggedIn() bool {
	return mc.loggedIn
}

func (mc MancalaClientProt) IsInLobby() bool {
	return mc.inLobby
}

func (mc *MancalaClientProt) LeaveLobby() {
	mc.inLobby = false
}

func debugLog(a ...interface{}) {
	log.Println(a...)
}

func (mc MancalaClientProt) printCookies() {
	sUrl, _ := url.Parse(mc.url)
	for _, v := range mc.client.Jar.Cookies(sUrl) {
		debugLog("Cookie:", v)
	}
}

func (mc *MancalaClientProt) SignUp(username, password string) error {
	username = url.QueryEscape(username)
	password = url.QueryEscape(password)
	urlStr := mc.url + "api/user/new"
	resp, err := mc.client.Post(urlStr, "application/x-www-form-urlencoded", strings.NewReader("username="+username+"&password="+password))
	if err != nil {
		return err
	}
	var content []byte
	var respError errorMessage
	content, err = ioutil.ReadAll(resp.Body)
	//debugLog(urlStr, username, password)
	//debugLog(string(content))
	if err != nil {
		return err
	}
	dec := json.NewDecoder(bytes.NewReader(content))
	dec.DisallowUnknownFields()
	err = dec.Decode(&respError)
	if err == nil {
		return errors.New(respError.Error)
	}
	mc.loggedIn = true
	return nil
}

func (mc *MancalaClientProt) SignIn(username, password string) error {
	username = url.QueryEscape(username)
	password = url.QueryEscape(password)
	resp, err := mc.client.Post(mc.url+"api/user/login", "application/x-www-form-urlencoded", strings.NewReader("username="+username+"&password="+password))
	if err != nil {
		return err
	}
	var content []byte
	var respError errorMessage
	content, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	dec := json.NewDecoder(bytes.NewReader(content))
	dec.DisallowUnknownFields()
	err = dec.Decode(&respError)
	if err == nil {
		return errors.New(respError.Error)
	}
	mc.loggedIn = true
	return nil
}

func (mc *MancalaClientProt) CreateLobby() (string, error) {
	urlStr := mc.url + "api/lobby/create"
	resp, err := mc.client.Post(urlStr, "application/x-www-form-urlencoded", nil)
	if err != nil {
		return "", err
	}
	var content []byte
	var respError errorMessage
	content, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	//debugLog(url)
	//debugLog(string(content))
	dec := json.NewDecoder(bytes.NewReader(content))
	dec.DisallowUnknownFields()
	err = dec.Decode(&respError)
	if err == nil {
		return "", errors.New(respError.Error)
	}
	mc.inLobby = true
	return string(content), nil
}

func (mc *MancalaClientProt) JoinLobby(room string) error {
	room = url.QueryEscape(room)
	resp, err := mc.client.Post(mc.url+"api/lobby/join", "application/x-www-form-urlencoded", strings.NewReader("room="+room))
	if err != nil {
		return err
	}
	var content []byte
	var respError errorMessage
	content, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	dec := json.NewDecoder(bytes.NewReader(content))
	dec.DisallowUnknownFields()
	err = dec.Decode(&respError)
	if err == nil {
		return errors.New(respError.Error)
	}
	mc.inLobby = true
	return nil
}

func (mc *MancalaClientProt) GetGameState() (*McBoard, error) {
	resp, err := mc.client.Get(mc.url + "api/lobby/state")
	if err != nil {
		return nil, err
	}
	var content []byte
	var respError errorMessage
	content, err = ioutil.ReadAll(resp.Body)
	//debugLog(string(content))
	if err != nil {
		return nil, err
	}
	//err = json.Unmarshal(content, &respError)
	dec := json.NewDecoder(bytes.NewReader(content))
	dec.DisallowUnknownFields()
	err = dec.Decode(&respError)
	if err == nil {
		return nil, errors.New(respError.Error)
	}

	var board McBoard
	err = json.Unmarshal(content, &board)
	if err != nil {
		return nil, err
	}
	return &board, nil
}

func (mc *MancalaClientProt) MakeTurn(cell int) (*McBoard, error) {
	resp, err := mc.client.Post(mc.url+"api/lobby/turn", "application/x-www-form-urlencoded", strings.NewReader("cell="+fmt.Sprint(cell)))
	if err != nil {
		return nil, err
	}
	var content []byte
	var respError errorMessage
	content, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	dec := json.NewDecoder(bytes.NewReader(content))
	dec.DisallowUnknownFields()
	err = dec.Decode(&respError)
	if err == nil {
		return nil, errors.New(respError.Error)
	}

	var board McBoard
	err = json.Unmarshal(content, &board)
	if err != nil {
		return nil, err
	}
	return &board, nil
}

func newMancalaClient(url string) MancalaClient {
	jar, err := cookiejar.New(nil) //TODO cookies
	panicCheck(err)
	return &MancalaClientProt{
		url: url,
		client: &http.Client{
			Transport:     nil,
			CheckRedirect: nil,
			Jar:           jar,
			Timeout:       0,
		},
	}
}

func (mc MancalaClientProt) debugCookie() {
	sUrl, err := url.Parse(mc.url)
	panicCheck(err)
	for _, c := range mc.client.Jar.Cookies(sUrl) {
		log.Println(c.Name, ":", c.Value)
	}
}
