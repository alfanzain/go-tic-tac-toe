package main

import (
	"encoding/json"
	"log"
)

type PlayerActionType string
type PlayerSide string
type GameState string

const (
	PlayerActionTypeInit PlayerActionType = "init"
	PlayerActionTypeMove PlayerActionType = "move"
	PlayerActionTypeEnd  PlayerActionType = "end"

	PlayerSideX PlayerSide = "X"
	PlayerSideO PlayerSide = "O"

	GameStateOngoing GameState = "ongoing"
	GameStateXWins   GameState = "X wins"
	GameStateOWins   GameState = "O wins"
	GameStateDraw    GameState = "draw"
)

type PlayerMovementPosition struct {
	Row int `json:"row"`
	Col int `json:"col"`
}

type PlayerActionPayload struct {
	Side     PlayerSide                           `json:"side,omitempty"`
	Position PlayerMovementPosition               `json:"position,omitempty"`
	Board    [BoardMaxRow][BoardMaxCol]PlayerSide `json:"board"`
	Status   GameState                            `json:"status"`
	Actor    *Player                              `json:"actor"`
}

type PlayerAction struct {
	ActionType PlayerActionType    `json:"action_type"`
	Data       PlayerActionPayload `json:"data"`
}

func NewPlayerAction(player *Player, p []byte) *PlayerAction {
	var playerAction PlayerAction

	log.Println("unmarshaling", string(p))
	err := json.Unmarshal(p, &playerAction)
	playerAction.Data.Actor = player
	playerAction.Data.Side = player.Side

	if err != nil {
		log.Fatal("failed unmarshal chat: ", err)
	}

	return &playerAction
}

func InitPlayerAction() *PlayerAction {
	return &PlayerAction{
		ActionType: PlayerActionTypeInit,
		Data: PlayerActionPayload{
			Board:  [3][3]PlayerSide{},
			Status: GameStateOngoing,
			Actor:  &Player{},
		},
	}
}

func EndPlayerAction(room *Room, player *Player) *PlayerAction {
	var gameState GameState
	if player.Side == PlayerSideX {
		gameState = GameStateOWins
	} else {
		gameState = GameStateXWins
	}
	return &PlayerAction{
		ActionType: PlayerActionTypeEnd,
		Data: PlayerActionPayload{
			Board:  room.Board,
			Status: gameState,
			Actor:  &Player{},
		},
	}
}

func ParsePlayerAction(v any) []byte {
	log.Println("marshaling", v)
	p, err := json.Marshal(v)
	if err != nil {
		log.Fatal("failed marshal message: ", err)
	}

	return p
}
