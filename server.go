package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/gorilla/websocket"
)

type Server struct {
	Rooms map[string]*Room
}

type CreateGameRoomResponse struct {
	Room *Room `json:"room"`
}

type DebugGameRoomListResponse struct {
	Rooms     map[string]*Room `json:"rooms"`
	RoomCount int              `json:"count"`
}

var upgrader = websocket.Upgrader{}

func NewServer() *Server {
	return &Server{
		Rooms: map[string]*Room{},
	}
}

func (s *Server) ServeWebMainMenu(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "scenes/index.html")
}

func (s *Server) ServeCreateGameRoom(w http.ResponseWriter, r *http.Request) {
	log.Println("creating a game room...")

	room := NewRoom()
	log.Println("room", room)

	log.Println("registering room...")
	s.Rooms[room.Id] = room

	log.Println("server rooms now: ", s.Rooms)

	go room.RunMessaging()

	render.Render(w, r, NewSuccessResponse(CreateGameRoomResponse{
		Room: room,
	}))
}

func (s *Server) ServeWebGameRoom(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		log.Println("id is required")
		http.ServeFile(w, r, "scenes/404.html")
		return
	}

	log.Println("accessing game room id:", id)

	if _, exists := s.Rooms[id]; !exists {
		log.Println("game room not found, id:", id)
		http.ServeFile(w, r, "scenes/404.html")
		return
	}

	http.ServeFile(w, r, "scenes/room.html")
}

func (s *Server) ServeGameWebsocket(w http.ResponseWriter, r *http.Request) {
	log.Println("upgrading http to websocket...")

	roomId := chi.URLParam(r, "id")
	if roomId == "" {
		log.Println("room id is required")
		return
	}

	room, ok := s.Rooms[roomId]
	log.Println("room", room)
	if !ok {
		log.Println("game room not found")
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("failed upgrade http to websocket", err)
		return
	}

	player := NewPlayer(room, conn)
	s.Rooms[roomId].register <- player

	go player.Read()
	go player.Write()
}

func (s *Server) ServeDebugGameRoomList(w http.ResponseWriter, r *http.Request) {
	render.Render(w, r, NewSuccessResponse(DebugGameRoomListResponse{
		Rooms:     s.Rooms,
		RoomCount: len(s.Rooms),
	}))
}
