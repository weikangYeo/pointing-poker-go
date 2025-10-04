package entity

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	// todo might need a proper client id, so in FE side when vote state is boardcasted, FE can map who is voted who is not.
	Name string
	Room *Room
	// websocket connection
	Conn *websocket.Conn
	// outbound message
	Send        chan SocketMessage
	CurrentVote string
}

const (
	maxMessageSize = 512
	// Time allowed to read the next pong message from the client.
	pongWait = 90 * time.Second
	pingWait = 60 * time.Second
)

// ReceiveMessageFromSocket Establish Connection from websocket and send to Room's Boardcast channel if any
// it can accept ctx context.Context too, given if this need to share a same context if this func need to call to other service or database call (and cancel it)
func (client *Client) ReceiveMessageFromSocket() {
	// to close connection after exit (either due to an error or a clean close)
	defer client.Conn.Close()

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
		log.Println("Received message:", string(message))

		var socketMessage SocketMessage
		if err := json.Unmarshal(message, &socketMessage); err != nil {
			log.Printf("Unexpected Message format error: %v", err)
			// todo, think of better logic flow
			continue
		}

		if socketMessage.Action == "vote" {
			var vote VoteReq
			if err := json.Unmarshal(socketMessage.Payload, &vote); err != nil {
				log.Printf("Unexpected Message format error: %v", err)
			}
			client.CurrentVote = vote.Point
			client.Room.BroadcastVoteState()
		}
	}
}

func (client *Client) SendMessage(ctx context.Context) {
	ticker := time.NewTicker(60 * time.Second)
	defer func() {
		ticker.Stop()
		client.Conn.Close()
	}()
	for {
		select {

		case <-ctx.Done():
			// this trigger when client disconnected -> trigger client.Conn.ReadMessage() error in ReceiveMessageFromSocket and return back ConnectToRoom
			// and triggered `defer cancel()` which will signal ctx.Done()
			log.Printf("SendMessage for client %s: context canceled", client.Name)
			client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		case message, ok := <-client.Send:
			log.Println("Sending message out to client")
			client.Conn.SetWriteDeadline(time.Now().Add(pingWait))
			if !ok {
				// with ctx.Done(), this might not come in anymore
				// hub closed channel
				log.Println("error !")
				client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			writer, err := client.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			jsonMsg, err := json.Marshal(message)
			if err != nil {
				log.Println("could not marshal json:", err)
				continue // Or handle error appropriately
			}
			writer.Write(jsonMsg)

			// write our message that still in send channel (if any)
			remainingMessageCount := len(client.Send)
			for i := 0; i < remainingMessageCount; i++ {
				nextMessage := <-client.Send
				jsonMsg, err = json.Marshal(nextMessage)
				if err != nil {
					log.Println("could not marshal json:", err)
					continue // Or handle error appropriately
				}
				writer.Write(jsonMsg)
			}
			e := writer.Close()
			if e != nil {
				return
			}
		case <-ticker.C:
			client.Conn.SetWriteDeadline(time.Now().Add(pingWait))
			if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}

	}
}
