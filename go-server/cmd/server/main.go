package main

import (
	"liar-of-turing/common"
	"liar-of-turing/internal/handlers"
	"liar-of-turing/services"
	"log"
	"net/http"

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

	common.SetFastAPIURL()

	log.Println("Starting channel listener ")
	userManager := services.NewUserManager()
	webSocketService := services.NewWebSocketService()
	gameState := services.NewGameState()
	mux := routes(userManager, webSocketService, gameState)
	go handlers.ListenToWebSocketChannel(userManager, webSocketService, gameState)
	// go handlers.ListenToGPTWebSocketChannel(FastAPIURL)

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

func routes(userManager *services.UserManager, webSocketService *services.WebSocketService, gameState *services.GameState) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleWebSocketRequest(w, r, userManager, webSocketService, gameState)
	})
	mux.HandleFunc("/withGPT", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleGPTWebSocketRequest(w, r, webSocketService)
	})

	return mux
}
