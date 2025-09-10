package handler

import (
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

	client := &entity.Client{
		Name: username,
		Conn: conn,
		Room: room,
		Send: make(chan []byte, 256),
	}

	client.Room.RegisterChan <- client

	go client.SendMessage()
	// dont go routine here as we need this become blocked and own this resource
	// so when this exit, defer triggered, and connection closed.
	client.ReceiveMessageFromSocket()
}
