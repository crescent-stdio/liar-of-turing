package main

import (
	"liarOfTuring/internal/handlers"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

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

	FastAPIURL := os.Getenv("FASTAPI_URL")

	mux := routes()
	log.Println("Starting channel listener ")
	go handlers.ListenToWsChannel()
	go handlers.ListenToGPTWsChannel(FastAPIURL)

	// bty env
	server := http.Server{
		Addr:    ":8443",
		Handler: c.Handler(mux),
	}

	log.Println("Listening on port :8443")
	err = server.ListenAndServe()

	if err != nil {
		log.Fatal("Error starting server: ", err)
	}

}

func routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", handlers.WsEndpoint)
	mux.HandleFunc("/withGPT", handlers.WithGPTWsEndpoint)

	return mux
}
