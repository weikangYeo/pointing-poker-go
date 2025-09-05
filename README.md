# Pointing Poker Go (Websocket)
This project is aim to try out websocket implementation in Golang by developing a Pointing Poker System.

## TODO
- [ ] Basic Setup of Hub Client Connection
- [ ] Sample Web to test out idea
- [ ] Logic to Vote
  - [ ] Do not accept input when isVisible = true
- [ ] Logic to Show All Vote
  - [ ] Do i need a temp storage to store vote? (if card is all hide, should not send message to clients)
- [ ] Logic to close Room after TTL
- [ ] Logic to control joiner in Room
- [ ] [Optional] Deploy to cluster to test multi client
