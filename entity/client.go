package entity

import (
	"bytes"
	"context"
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
		client.Room.BroadcastChan <- message
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
			client.Conn.SetWriteDeadline(time.Now().Add(pingWait))
			if !ok {
				// with ctx.Done(), this might not come in anymore
				// hub closed channel
				client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			writer, err := client.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			writer.Write(message)

			// write our message that still in send channel
			remainingMessageCount := len(client.Send)
			for i := 0; i < remainingMessageCount; i++ {
				writer.Write([]byte{10})
				writer.Write(<-client.Send)
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
