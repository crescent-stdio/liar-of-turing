package main

import (
	"liarOfTuring/internal/handlers"
	"net/http"

	"github.com/bmizerany/pat"
)

// routes defines theapplication routes
func routes() http.Handler {
	mux := pat.New()

	// making a chat room with nextjs
	mux.Get("/", http.HandlerFunc(handlers.Home))
	mux.Get("/ws", http.HandlerFunc(handlers.WsEndpoint))
	// mux.Get("/api/room", http.HandlerFunc(handlers.CreateRoom))
	// mux.Get("/api/room/:roomID", http.HandlerFunc(handlers.GetRoomInfo))

	// r.HandleFunc("/send", SendMessage)
	// r.HandleFunc("/join", JoinRoom)
	// r.HandleFunc("/leave", LeaveRoom)
	// r.HandleFunc("/rooms", GetRooms)
	// r.HandleFunc("/room/{roomID}", GetRoomInfo)
	// r.HandleFunc("/users", GetUsers)
	// r.HandleFunc("/user/{username}", GetUserInfo)
	// r.HandleFunc("/user/{username}/rooms", GetUserRooms)
	// r.HandleFunc("/user/{username}/room/{roomID}", GetUserRoomInfo)
	// r.HandleFunc("/user/{username}/room/{roomID}/messages", GetUserRoomMessages)
	// r.HandleFunc("/user/{username}/room/{roomID}/message/{messageID}", GetUserRoomMessageInfo)

	return mux
}
