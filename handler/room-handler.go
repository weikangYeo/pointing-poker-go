package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"pointing-poker-go/entity"
	"sync"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type RoomServer struct {
	// like `sync` claus in java, to lock object for safe concurrency purpose
	mu    sync.RWMutex
	rooms map[string]*entity.Room
}

func NewRoomHandler() *RoomServer {
	return &RoomServer{
		rooms: make(map[string]*entity.Room),
	}
}

func (server *RoomServer) CreateRoom(w http.ResponseWriter, r *http.Request) {
	room := entity.NewRoom()
	if room == nil {
		log.Println("CreateRoom Error")
		writeHttpResponse(w, http.StatusInternalServerError, map[string]string{
			"message": "Create Room Error",
		})
		return
	}

	server.mu.Lock()
	server.rooms[room.Id] = room
	go room.Start()
	server.mu.Unlock()

	writeHttpResponse(w, http.StatusCreated, map[string]string{
		"id": room.Id,
	})
}

func (server *RoomServer) ConnectToRoom(w http.ResponseWriter, r *http.Request) {
	roomId := mux.Vars(r)["id"]
	username := r.URL.Query().Get("username")
	if username == "" {
		writeHttpResponse(w, http.StatusBadRequest, map[string]string{
			"message": "Username can not be empty",
		})
	}

	server.mu.RLock()
	room, ok := server.rooms[roomId]
	server.mu.RUnlock()

	if !ok {
		writeHttpResponse(w, http.StatusNotFound, map[string]string{
			"message": fmt.Sprintf("Connect To Room Error, room id %s not found", roomId),
		})
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading ws: %v", err)
		return
	}

	// Create a new context for this client connection
	// it's child of the request context, so it will also be cancelled when the http connection is closed
	clientCtx, cancel := context.WithCancel(r.Context())
	// it ensured that when ReceiveMessageFromSocket return (due to client disconnected) the context is canceled.
	defer cancel()

	client := &entity.Client{
		Name: username,
		Conn: conn,
		Room: room,
		Send: make(chan entity.SocketMessage),
	}

	client.Room.RegisterChan <- client

	go client.SendMessage(clientCtx)
	//log.Println("Update status all connected clients")
	//room.BroadcastVoteState()
	// dont go routine here as we need this become blocked and own this resource
	// so when this exit, defer triggered, and connection closed.
	client.ReceiveMessageFromSocket()

	// after client is disconnected (client.ReceiveMessageFromSocket returned), unregister from room
	client.Room.UnregisterChan <- client
}
