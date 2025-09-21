# Pointing Poker Go (Websocket)

This project is aim to try out websocket implementation in Golang by developing a Pointing Poker System.

## Functional Overview

- This repo is act as a BE repository of pointing poker app.
- A Client can connect to a Room at a time.
- A Room will be destroyed when no client is connected to it

## References

### Ping Pong Health Check Flow

1. Server Sends Ping: The server's writePump sends a PingMessage to the client, effectively asking, "Are you still
   there?"
2. Client Responds with Pong: The browser's native WebSocket implementation sees the Ping and automatically sends a
   PongMessage back. You don't need to write any JavaScript for this; it's part of the
   standard.
3. Server Waits for Pong: The server's readPump has a read deadline set (pongWait). When the PongMessage arrives from
   the client, the PongHandler is triggered, and it resets the deadline.
4. The "Dead Client" Scenario: If the client has disconnected (e.g., the user closed the browser tab, lost internet), it
   won't receive the server's Ping. Therefore, it will never send a Pong back. The
   server's read deadline will expire, c.conn.ReadMessage() will return a timeout error, and the readPump will exit,
   closing the connection and cleaning up resources.

### Communication Flow

1. FE send data message, `conn = new WebSocket("ws://" + document.location.host + "/ws");`
2. goroutine client.go ReceiveMessageFromSocket() blocked at `_, message, err := client.Conn.ReadMessage()`
3. ReceiveMessageFromSocket() send message to Room's boardcast channel
4. Room's goroutine Start() waited case, `select case message := <-room.BroadcastChan:` run, and send to each Client's
   `send` channel
5. Client's WriteMessageToSocket select case run, which listend to `send` channel, and TBC...

## TODO

- [ ] Basic Setup of Hub Client Connection
- [ ] Sample Web to test out idea
- [ ] Logic to Vote
    - [ ] Do not accept input when isVisible = true
    - [ ] boardcast isvoted logic, room might need to keep isvoted stated, or shall room keep vote too?
    - [ ] when join room, get room latest state
- [ ] Logic to Show All Vote
    - [ ] Do i need a temp storage to store vote? (if card is all hide, should not send message to clients)
- [ ] Logic to close Room after TTL
- [ ] Logic to control joiner in Room
- [ ] Logic to control same user can't appear twice  
- [ ] [Optional] Deploy to cluster to test multi client

