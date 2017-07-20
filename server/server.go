package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var (
	MESSAGE    = "MESSAGE"
	CONNECT    = "CONNECT"
	DISCONNECT = "DISCONNECT"
)

var users []user

type user struct {
	username string
	conn     *websocket.Conn
}

type message struct {
	Command   string
	Timestamp time.Time
	Text      string
	Username  string
}

func removeUser(user user) {
	for index, existing := range users {
		if existing.username == user.username {
			users = append(users[:index], users[index+1:]...)
			return
		}
	}
}

func parseCommands(user user) {
	for true {
		var msg message
		if err := user.conn.ReadJSON(&msg); err != nil {
			fmt.Printf("%v disconnected ungracefully\n", user.username)
			removeUser(user)
			return
		}
		fmt.Printf("%#v\n", msg)

		switch msg.Command {
		case MESSAGE:
			broadcast(msg)
		case DISCONNECT:
			fmt.Printf("%v disconnected\n", user.username)
			removeUser(user)
			return
		}
	}
}

func broadcast(message message) {
	for _, user := range users {
		if user.username == message.Username {
			continue
		}
		if err := user.conn.WriteJSON(message); err != nil {
			panic(err)
		}
	}
}

func connectHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		panic(err)
	}

	var msg message
	if err := conn.ReadJSON(&msg); err != nil {
		panic(err)
	}

	var user user
	user.username = msg.Username
	user.conn = conn

	users = append(users, user)

	go parseCommands(user)

	fmt.Printf("%v connected\n", user.username)
}

func main() {
	port := flag.String("Port", "12345", "Port to bind server to.")

	http.HandleFunc("/", connectHandler)
	if err := http.ListenAndServe(":"+*port, nil); err != nil {
		panic("Error" + err.Error())
	}
}
