package entity

type Room struct {
	id          string
	roundId     int
	showAllCard bool
	// todo might consider use k = name, v = client
	joinedClients map[*Client]bool
	// inbound message from clients publish here
	broadcastChan chan []byte
	// register request publish here
	registerChan chan *Client
	// unregister request publish here
	unregisterChan chan *Client
}
