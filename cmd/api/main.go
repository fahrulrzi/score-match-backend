package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fahrulrzi/score-match-backend/configs"
	"github.com/fahrulrzi/score-match-backend/internal/delivery/http/handler"
	"github.com/fahrulrzi/score-match-backend/internal/delivery/http/middleware"
	"github.com/fahrulrzi/score-match-backend/internal/delivery/http/routes"
	"github.com/fahrulrzi/score-match-backend/internal/infrastructure/database"
	"github.com/fahrulrzi/score-match-backend/internal/repository/postgres"
	"github.com/fahrulrzi/score-match-backend/internal/usecase"
	"github.com/fahrulrzi/score-match-backend/pkg/jwt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
)

func main() {
	config, err := configs.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Initialize database connection
	db, err := database.NewPostgresConnection(&config.Database)
	if err != nil {
		log.Fatal("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize database schema
	err = database.InitDatabase(db)
	if err != nil {
		log.Fatal("Failed to initialize database schema: %v", err)
	}

	// Initialize repositories
	tokenRepo := postgres.NewTokenRepository(db)
	userRepo := postgres.NewUserRepository(db)
	customerRepo := postgres.NewCustomerRepository(db)

	// Initialize services
	jwtService := jwt.NewJWTService(config.JWT.Secret, config.JWT.ExpireTime)

	// Initialize use cases
	authUseCase := usecase.NewAuthUseCase(userRepo, jwtService)
	customerUseCase := usecase.NewCustomerUseCase(customerRepo)

	// Initialize HTTP handlers
	authHandler := handler.NewAuthHandler(authUseCase, jwtService)
	customerHandler := handler.NewCustomerHandler(customerUseCase)
	authMiddleware := middleware.NewAuthMiddleware(jwtService, tokenRepo)

	// Initialize router
	router := mux.NewRouter()
	routes.SetupRoutes(router, authMiddleware, authHandler, customerHandler)

	// Setup CORS middleware
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // frontend lo
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	// Wrap router with CORS
	handlerWithCORS := corsMiddleware.Handler(router)

	// Configure HTTP server
	server := &http.Server{
		Addr:         ":" + config.Server.Port,
		Handler:      handlerWithCORS,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Start HTTP server
	go func() {
		log.Printf("Starting server on port %s", config.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Printf("Shutting down server...")

	// Create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Doesn't block if no connection, but will otherwise wait until the timeout
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}
