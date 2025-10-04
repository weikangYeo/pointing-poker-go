package entity

import "encoding/json"

type VoteReq struct {
	Point string
}

type RoomVoteState struct {

	// todo change to client id -> boolean map
	IsVotedByClientNameMap map[string]bool
}

type SocketMessage struct {
	Action  string          `json:"action"`
	Payload json.RawMessage `json:"payload"`
}
