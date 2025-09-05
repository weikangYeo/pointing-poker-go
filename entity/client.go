package entity

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	name string
	room *Room
	// websocket connection
	conn *websocket.Conn
	// outbound message
	send chan []byte
}
