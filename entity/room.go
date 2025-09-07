package entity

import (
	"crypto/rand"
	"log"
)

type Room struct {
	Id          string
	RoundId     int
	ShowAllCard bool
	// todo might consider use k = name, v = client
	JoinedClients map[*Client]bool
	// inbound message from clients publish here
	BroadcastChan chan []byte
	// register request publish here
	RegisterChan chan *Client
	// unregister request publish here
	UnregisterChan chan *Client
}

func NewRoom() *Room {
	b := make([]byte, 12)
	_, err := rand.Read(b)
	if err != nil {
		log.Printf("Error generating random number: %v", err)
		return nil
	}
	roomId := string(b)

	return &Room{
		Id:             roomId,
		RoundId:        1,
		ShowAllCard:    false,
		JoinedClients:  make(map[*Client]bool),
		BroadcastChan:  make(chan []byte),
		RegisterChan:   make(chan *Client),
		UnregisterChan: make(chan *Client),
	}
}

func (room *Room) Start() {
	for {
		select {
		// register client
		case client := <-room.RegisterChan:
			room.JoinedClients[client] = true
		// unregister client
		case client := <-room.UnregisterChan:
			// mean if room.JoinedClients[client] exists, and assign to "ok" variable
			if _, ok := room.JoinedClients[client]; ok {
				delete(room.JoinedClients, client)
				// close channel, so client.SendMessage() can detect it and close websocket connection
				close(client.Send)
			}
		// boardcast message
		case message := <-room.BroadcastChan:
			for client := range room.JoinedClients {
				select {
				case client.Send <- message:
				// with default case, this "select" clause become not blocking
				// when a message send to client.send, if the client.send is fulled,
				// we treat it as client isn't reading message (i.e. disconnected)
				// so we close the connection, else this message := <-h.broadcast will be blocked for this client.
				default:
					close(client.Send)
					delete(room.JoinedClients, client)
				}

			}
		}
	}
}
