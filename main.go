package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"ws-go/handlers"
	"ws-go/server"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("error loading env file")
	}

	PORT := os.Getenv("PORT")
	JWT_SECRET := os.Getenv("JWT_SECRET")
	DATABASE_URL := os.Getenv("DATABASE_URL")

	s, err := server.NewServer(context.Background(), server.Config{
		Port:        PORT,
		JWTSecret:   JWT_SECRET,
		DatabaseURL: DATABASE_URL,
	})

	s.Start(BindRoutes)

}
func BindRoutes(s server.Server, r *mux.Router) {

	r.HandleFunc("/ws", s.Hub().HandleWebsocket)

	// r.Use(middleware.CheckAuthMiddleware(s))

	r.HandleFunc("/", handlers.HomeHandler(s)).Methods(http.MethodGet)
	r.HandleFunc("/signup", handlers.SignUpHandler(s)).Methods(http.MethodPost)
	r.HandleFunc("/login", handlers.LoginHandler(s)).Methods(http.MethodPost)
	r.HandleFunc("/me", handlers.MeHandler(s)).Methods(http.MethodGet)
}
