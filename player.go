package main

import (
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/goombaio/namegenerator"
	"github.com/gorilla/websocket"
)

type Player struct {
	Id   string          `json:"id"`
	Name string          `json:"name"`
	Room *Room           `json:"-"`
	Conn *websocket.Conn `json:"-"`

	send     chan *PlayerAction
	sendRoom chan *RoomMessage
}

func NewPlayer(room *Room, conn *websocket.Conn) *Player {
	seed := time.Now().UTC().UnixNano()
	nameGenerator := namegenerator.NewNameGenerator(seed)
	name := nameGenerator.Generate()

	return &Player{
		Id:   uuid.NewString(),
		Name: name,
		Room: room,
		Conn: conn,

		send:     make(chan *PlayerAction),
		sendRoom: make(chan *RoomMessage),
	}
}

func (player *Player) Read() {
	log.Printf("player [%s:%s] read() is running\n", player.Id, player.Name)

	defer func() {
		player.Room.unregister <- player
		player.Conn.Close()
	}()

	for {
		_, p, err := player.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("failed read chat: %v", err)
			}
			return
		}
		if p == nil {
			log.Println("payload is nil")
			continue
		}

		playerAction := NewPlayerAction(player, p)
		log.Printf("player [%s:%s] makes an action: %v\n", player.Id, player.Name, playerAction)

		log.Println("player", player.Name, "is broadcasting the message")
		player.Room.broadcast <- playerAction
	}
}

func (player *Player) Write() {
	log.Printf("player [%s:%s] write() is running\n", player.Id, player.Name)

	for {
		select {
		case action, ok := <-player.send:
			if !ok {
				player.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			log.Printf("player [%s:%s] is about to receiving an action from its opponent: %v\n", player.Id, player.Name, action)

			writer, err := player.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Fatal("failed to initiate next writer", err)
				return
			}

			log.Println("player", player.Name, "is receiving the action")
			writer.Write(ParsePlayerAction(action))

			if err = writer.Close(); err != nil {
				log.Fatal("failed to initiate next writer", err)
				return
			}
		case message, ok := <-player.sendRoom:
			if !ok {
				player.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			log.Printf("player [%s:%s] is about to receiving a message from server: %v\n", player.Id, player.Name, message)

			writer, err := player.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Fatal("failed to initiate next writer", err)
				return
			}

			log.Println("player", player.Name, "is receiving the message")
			writer.Write(message.Parse())

			if err = writer.Close(); err != nil {
				log.Fatal("failed to initiate next writer", err)
				return
			}
		}
	}
}
