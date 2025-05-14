package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/goombaio/namegenerator"
)

type RoomState string

const (
	BoardMaxRow int = 3
	BoardMaxCol int = 3

	RoomStateWaiting  RoomState = "waiting"
	RoomStateActive   RoomState = "active"
	RoomStateFinished RoomState = "finished"
)

type Room struct {
	Id          string                               `json:"id"`
	Players     [2]*Player                           `json:"players"`
	Board       [BoardMaxRow][BoardMaxCol]PlayerSide `json:"board"`
	CurrentTurn PlayerSide                           `json:"current_turn"`
	Status      RoomState                            `json:"status"`

	register      chan *Player
	unregister    chan *Player
	broadcast     chan *PlayerAction
	broadcastRoom chan *RoomMessage
}

func NewRoom() *Room {
	seed := time.Now().UTC().UnixNano()
	nameGenerator := namegenerator.NewNameGenerator(seed)
	id := nameGenerator.Generate()

	return &Room{
		Id:          id,
		Players:     [2]*Player{},
		Board:       [3][3]PlayerSide{},
		CurrentTurn: PlayerSideX,
		Status:      RoomStateWaiting,

		register:      make(chan *Player),
		unregister:    make(chan *Player),
		broadcast:     make(chan *PlayerAction),
		broadcastRoom: make(chan *RoomMessage),
	}
}

func (room *Room) IsFull() bool {
	return room.Players[0] != nil && room.Players[1] != nil
}

func (room *Room) RunMessaging() {
	log.Println("Game room is waiting for players...")

	for {
		select {
		case player := <-room.register:
			log.Printf("player [%s:%s] joining game [%s]...\n", player.Id, player.Name, room.Id)

			if room.Players[0] != nil && room.Players[1] != nil {
				log.Printf("player [%s:%s] wants to join game [%s], turns out it's full\n", player.Id, player.Name, room.Id)
				continue
			}

			if room.Players[0] == nil {
				room.Players[0] = player
				log.Printf("player [%s:%s] joined game [%s] as X!\n", player.Id, player.Name, room.Id)
			} else if room.Players[1] == nil {
				room.Players[1] = player
				log.Printf("player [%s:%s] joined game [%s] as O!\n", player.Id, player.Name, room.Id)
			}

			for _, roomPlayer := range room.Players {
				if roomPlayer != nil {
					roomPlayer.sendRoom <- room.Info(roomPlayer)
				}
			}
		case message := <-room.broadcastRoom:
			log.Printf("broadcast room message to all players in the room [%s]\n", room.Id)

			for _, player := range room.Players {
				if player != nil {
					player.sendRoom <- message
				}
			}
		}
	}
}

func (room *Room) StartGame() {
	log.Printf("Game room [%s] started!\n", room.Id)

	for {
		select {
		case action := <-room.broadcast:
			log.Printf("broadcast action to all players in the room [%s]\n", room.Id)

			for _, player := range room.Players {
				player.send <- action
			}
		}
	}
}

func (room *Room) CheckGameStatus() GameState {
	winCoords := [8][3][2]int{
		{{0, 0}, {1, 0}, {2, 0}},
		{{0, 1}, {1, 1}, {2, 1}},
		{{0, 2}, {1, 2}, {2, 2}},

		{{0, 0}, {0, 1}, {0, 2}},
		{{1, 0}, {1, 1}, {1, 2}},
		{{2, 0}, {2, 1}, {2, 2}},

		{{0, 0}, {1, 1}, {2, 2}},
		{{0, 2}, {1, 1}, {2, 0}},
	}

	// X player
	for _, coords := range winCoords {
		if room.Board[coords[0][0]][coords[0][1]] == PlayerSideX && room.Board[coords[1][0]][coords[1][1]] == PlayerSideX && room.Board[coords[2][0]][coords[2][1]] == PlayerSideX {
			return GameStateXWins
		}
	}

	// O player
	for _, coords := range winCoords {
		if room.Board[coords[0][0]][coords[0][1]] == PlayerSideO && room.Board[coords[1][0]][coords[1][1]] == PlayerSideO && room.Board[coords[2][0]][coords[2][1]] == PlayerSideO {
			return GameStateOWins
		}
	}

	return GameStateOngoing
}

type RoomMessage struct {
	MessageType string `json:"message_type"`
	Data        any    `json:"data"`
}

func (msg *RoomMessage) Parse() []byte {
	log.Println("marshaling", msg)
	p, err := json.Marshal(msg)
	if err != nil {
		log.Fatal("failed marshal message: ", err)
	}

	return p
}

type GameInfoPayload struct {
	Room   *Room   `json:"room"`
	Player *Player `json:"player"`
}

func (room *Room) Info(player *Player) *RoomMessage {
	return &RoomMessage{
		MessageType: "game_info",
		Data: &GameInfoPayload{
			Room:   room,
			Player: player,
		},
	}
}
