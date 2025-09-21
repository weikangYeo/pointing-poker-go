package main

import (
	"fmt"
	"net/http"
	"pointing-poker-go/handler"

	"github.com/gorilla/mux"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {

	roomServer := handler.NewRoomHandler()

	router := mux.NewRouter()
	router.HandleFunc("/rooms", roomServer.CreateRoom).Methods("POST")
	router.HandleFunc("/rooms/{id}", roomServer.ConnectToRoom).Methods("GET")
	router.PathPrefix("/").Handler(http.FileServer(http.Dir(".")))
	//// todo
	//router.HandleFunc("/rooms/{id}/show", handler.ConnectToRoom).Methods("PATCH")
	//// todo
	//router.HandleFunc("/rooms/{id}/hide", handler.ConnectToRoom).Methods("PATCH")
	//// todo
	//router.HandleFunc("/rooms/{id}/votes", handler.ConnectToRoom).Methods("POST")
	fmt.Printf("Listening to Port 8080\n")
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		fmt.Println(err)
	}
}
