# Pointing Poker Go (Websocket)

This project is aim to try out websocket implementation in Golang by developing a Pointing Poker System.

## Functional Overview

- This repo is act as a BE repository of pointing poker app. 
- A Client can connect to a Room at a time.
- A Room will be destroyed when no client is connected to it

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
- [ ] [Optional] Deploy to cluster to test multi client

