package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	server := NewServer()

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", server.ServeWebMainMenu)
	r.Post("/room", server.ServeCreateGameRoom)
	r.Get("/room/{id:[a-z-]+}", server.ServeWebGameRoom)
	r.Get("/room/{id:[a-z-]+}/socket", server.ServeGameWebsocket)

	r.Get("/room", server.ServeDebugGameRoomList)

	log.Println("running server at :3000")
	http.ListenAndServe(":3000", r)
}
