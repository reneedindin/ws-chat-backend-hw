package client

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"golang.org/x/net/websocket"
)

type Message struct {
	Message string `json:"message"`
}

func NewClient(host, port string) {
	ws, err := connect(host, port)
	if ws != nil {
		defer ws.Close()
	}

	if err != nil {
		log.Fatalf("Connect to websoket server failed. err: %v\n", err)
	}

	go receive(ws)
	send(ws)
}

func connect(host, port string)(*websocket.Conn, error) {
	origin := "http://localhost/"
	url := fmt.Sprintf("ws://%s:%s/ws", host, port)
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		return nil, err
	}

	return ws, nil
}

func receive(ws *websocket.Conn) {
	var m Message

	for{
		if err := websocket.JSON.Receive(ws, &m); err != nil {
			log.Fatalf("Receive message is error, err: %s\n", err.Error())
		}

		fmt.Printf("Received: %s\n", m.Message)
	}
}

func send(ws *websocket.Conn) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		if text == "" {
			continue
		}

		m := Message{
			Message: text,
		}

		err := websocket.JSON.Send(ws, m)
		if err != nil {
			fmt.Printf("Send message is error, err: %s\n", err.Error())
			break
		}
	}
}