package routes

import (
	"net/http"

	"github.com/fahrulrzi/score-match-backend/internal/delivery/http/handler"
	"github.com/gorilla/mux"
)

func SetupRoutes(
	router *mux.Router,
	authHandler *handler.AuthHandler,
	customerHandler *handler.CustomerHandler,
) {
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	// Auth
	router.HandleFunc("/api/auth/register", authHandler.Register).Methods("POST")
	router.HandleFunc("/api/auth/login", authHandler.Login).Methods("POST")

	// Customer
	router.HandleFunc("/api/customer/score", customerHandler.GetFinalScore).Methods("POST")
}
