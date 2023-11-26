package main

import (
	"liarOfTuring/internal/handlers"
	"log"
	"net/http"
	"sync"

	"github.com/rs/cors"
)

// Counter stores the value of the counter
type Counter struct {
	Value int `json:"value"`
	mu    sync.Mutex
}

type Message struct {
	UserName  string `json:"username"`
	ID        string `json:"id"`
	Role      string `json:"role"`
	Timestamp int64  `json:"timestamp"`
}

type User struct {
	UserName string `json:"username"`
	Role     string `json:"role"`
}

type Room struct {
	RoomID string `json:"roomID"`
	Users  []User `json:"users"`
}

type RoomList struct {
	Rooms []Room `json:"rooms"`
}

type RoomInfo struct {
	RoomID string `json:"roomID"`
	Users  []User `json:"users"`
}

// port is the port to listen on
const port = ":8080"

// main is the main function
func main() {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},                                       // 모든 도메인에서의 요청 허용
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, // 허용할 메소드
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           3000, // 5분
	})

	mux := routes()
	log.Println("Starting channel listener")
	go handlers.ListenToWsChannel()

	log.Println("Starting server on port", port)

	_ = http.ListenAndServe(port, c.Handler(mux))
}
