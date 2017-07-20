package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"time"
)

var conn *websocket.Conn

var (
	MESSAGE    = "MESSAGE"
	CONNECT    = "CONNECT"
	DISCONNECT = "DISCONNECT"
)

type message struct {
	Command   string
	Timestamp time.Time
	Text      string
	Username  string
}

func connect(url, username string) (*websocket.Conn, error) {
	var dialer *websocket.Dialer
	conn, _, err := dialer.Dial(url, nil)
	if err != nil {
		return nil, err
	}
	conn.WriteJSON(&message{
		CONNECT,
		time.Now(),
		"",
		username,
	})

	return conn, nil
}

func send(msg message) {
	if err := conn.WriteJSON(msg); err != nil {
		panic(err)
	}
}

func receive(conn *websocket.Conn) {
	for true {
		var msg message
		if err := conn.ReadJSON(&msg); err != nil {
			panic(err)
		}
		fmt.Printf("%v:%v\n", msg.Username, msg.Text)
	}
}

func main() {
	url := "ws://localhost:12345/"
	var username string
	flag.StringVar(&username, "username", "", "Username")
	flag.Parse()
	if username == "" {
		fmt.Println("username required")
		return
	}

	var err error
	conn, err = connect(url, username)
	if err != nil {
		panic(err)
	}

	go receive(conn)

	for true {
		fmt.Print("> ")
		var text string
		fmt.Scanln(&text)
		var msg message
		msg.Username = username
		msg.Timestamp = time.Now()
		if text == "exit" {
			msg.Command = DISCONNECT
			send(msg)
			return
		} else {
			msg.Text = text
			msg.Command = MESSAGE
		}
		send(msg)
	}
}
