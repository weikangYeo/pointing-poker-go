package entity

import (
	"bytes"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	Name string
	Room *Room
	// websocket connection
	Conn *websocket.Conn
	// outbound message
	Send        chan []byte
	CurrentVote int
}

const (
	maxMessageSize = 512
	// Time allowed to read the next pong message from the client.
	pongWait = 90 * time.Second
)

// ReceiveMessageFromSocket Establish Connection from websocket and send to Room's Boardcast channel if any
func (client *Client) ReceiveMessageFromSocket() {
	// to close connection after exit (either due to an error or a clean close
	defer func() {
		client.Room.UnregisterChan <- client
		client.Conn.Close()
	}()

	client.Conn.SetReadLimit(maxMessageSize)
	// set a keep alive check with FE client
	client.Conn.SetReadDeadline(time.Now().Add(pongWait))
	// What to do when a Pong received from FE client
	client.Conn.SetPongHandler(func(string) error {
		// refresh deadline after receive Pong
		client.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, message, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(message)
		// boardcast message here
		// todo probably I wont need this, change to parse message logic here
		// might... need to boardcast, but boardcast as voted.
		client.Room.BroadcastChan <- message
	}
}

// todo writeMessageToSocket()
