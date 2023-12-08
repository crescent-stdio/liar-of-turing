package main

import (
	"liarOfTuring/internal/handlers"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/joho/godotenv"
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

// main is the main function
func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	log.Println("Hello World")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // All origins
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           3000, // 5 min
	})

	mux := routes()
	log.Println("Starting channel listener")
	go handlers.ListenToWsChannel()

	// URL := os.Getenv("URL")

	URL := ":" + os.Getenv("PORT")
	log.Println("URL", URL)

	server := http.Server{
		Addr:    URL,
		Handler: c.Handler(mux),
	}

	// SSL/TLS certificate
	// err = server.ListenAndServe() // for local testing
	err = server.ListenAndServeTLS("./cert.pem", "./key.pem") // for production

	if err != nil {
		log.Fatal("ListenAndServeTLS: ", err)
	}

}

func routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", handlers.WsEndpoint)

	return mux
}