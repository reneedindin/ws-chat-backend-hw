package main

import (
	"fmt"
	"log"
	"net/http"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/net/websocket"
)

type client struct {
	userID string
	conn *websocket.Conn
}

type Message struct {
	Message string `json:"message"`
}

type SendMessage struct {
	Message
	userID string
}

type clientManager struct {
	clients map[string]*client
	registerClient chan *client
	unregisterClient chan *client
	sendMessage chan SendMessage
}

func newclientManager() *clientManager {
	return &clientManager{
		clients:          make(map[string]*client),
		registerClient:   make(chan *client),
		unregisterClient: make(chan *client),
		sendMessage:      make(chan SendMessage),
	}
}

func (cm *clientManager) manager() {
	for {
		select {
		case client := <-cm.registerClient:
			cm.register(client)
		case client := <-cm.unregisterClient:
			cm.unregister(client)
		case message := <-cm.sendMessage:
			cm.send(message)
		}
	}
}

func (cm *clientManager) register(c *client) {
	cm.clients[c.userID] = c
}

func (cm *clientManager) unregister(c *client) {
	delete(cm.clients, c.userID)
}

func (cm *clientManager) send(m SendMessage) {
	for _, client := range cm.clients {
		if m.userID == client.userID {
			continue
		}

		if err := websocket.JSON.Send(client.conn, m.Message); err != nil {
			log.Printf("Send message is error, userID: %s, err: %s\n", client.userID, err.Error())
			continue
		}
	}
}

func wsHandler(ws *websocket.Conn, clientManager *clientManager) {
	uuid, err := uuid.NewV4()
	if err != nil {
		log.Printf("New uuid is error, err: %s\n", err.Error())
		return
	}

	client := &client{
		userID: uuid.String(),
		conn: ws,
	}
	clientManager.registerClient <- client

	for {
		m := SendMessage{
			Message: Message{},
			userID: client.userID,
		}

		if err := websocket.JSON.Receive(ws, &m.Message); err != nil {
			clientManager.unregisterClient <- client
			log.Printf("Receive message is error, userID: %s, err: %s\n", client.userID, err.Error())
			break
		}
		log.Printf("Receive message: %s\n", m.Message)

		clientManager.sendMessage <- m
	}
}



func main() {
	clientManager := newclientManager()
	go clientManager.manager()

	fmt.Println("Start Websoket server.")
	http.Handle("/ws", websocket.Handler(func (ws *websocket.Conn) {
		wsHandler(ws, clientManager)
	}))
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}