package main

import (
	"encoding/json"
	"fmt"
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
	name := nameGenerator.Generate()
	id := fmt.Sprintf("room-%s", name)

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
				player.Side = PlayerSideX
				room.Players[0] = player
				log.Printf("player [%s:%s] joined game [%s] as X!\n", player.Id, player.Name, room.Id)
			} else if room.Players[1] == nil {
				player.Side = PlayerSideO
				room.Players[1] = player
				log.Printf("player [%s:%s] joined game [%s] as O!\n", player.Id, player.Name, room.Id)
			}

			for _, roomPlayer := range room.Players {
				if roomPlayer != nil {
					roomPlayer.sendRoom <- room.Info(roomPlayer)
				}
			}

			if room.IsFull() {
				log.Printf("game room [%s] is full, starting game...", room.Id)
				room.Status = RoomStateActive
				room.CurrentTurn = PlayerSideX

				go room.StartGame()

				time.Sleep(500 * time.Millisecond)

				room.broadcast <- InitPlayerAction()
			}
		case player := <-room.unregister:
			if player.Side == PlayerSideX {
				room.Players[0] = nil
			} else {
				room.Players[1] = nil
			}

			room.Status = RoomStateFinished
			room.broadcast <- EndPlayerAction(room, player)
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

			if action.ActionType == PlayerActionTypeMove {
				log.Printf("updating board room [%s]...\n", room.Id)
				room.UpdateBoard(action.Data.Position.Row, action.Data.Position.Col, action.Data.Side)
				action.Data.Board = room.Board
				action.Data.Status = room.CheckGameStatus()
				log.Printf("board room now [%v]\n", room.Board)
			}

			if action.Data.Status == GameStateXWins || action.Data.Status == GameStateOWins {
				action.ActionType = PlayerActionTypeEnd
			}

			for _, player := range room.Players {
				if player != nil {
					player.send <- action
				}
			}
		}
	}
}

func (room *Room) UpdateBoard(row, col int, side PlayerSide) {
	room.Board[row][col] = side
}

func (room *Room) CheckGameStatus() GameState {
	winCoords := [8][3][2]int{
		{{0, 0}, {0, 1}, {0, 2}}, // Row 0
		{{1, 0}, {1, 1}, {1, 2}}, // Row 1
		{{2, 0}, {2, 1}, {2, 2}}, // Row 2
		{{0, 0}, {1, 0}, {2, 0}}, // Column 0
		{{0, 1}, {1, 1}, {2, 1}}, // Column 1
		{{0, 2}, {1, 2}, {2, 2}}, // Column 2
		{{0, 0}, {1, 1}, {2, 2}}, // Diagonal top-left to bottom-right
		{{0, 2}, {1, 1}, {2, 0}}, // Diagonal top-right to bottom-left
	}

	for _, coords := range winCoords {
		a, b, c := coords[0], coords[1], coords[2]
		if room.Board[a[0]][a[1]] == PlayerSideX &&
			room.Board[b[0]][b[1]] == PlayerSideX &&
			room.Board[c[0]][c[1]] == PlayerSideX {
			return GameStateXWins
		}

		if room.Board[a[0]][a[1]] == PlayerSideO &&
			room.Board[b[0]][b[1]] == PlayerSideO &&
			room.Board[c[0]][c[1]] == PlayerSideO {
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
