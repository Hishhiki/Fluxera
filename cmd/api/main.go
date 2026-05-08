package main

import (
	"fluxera/internal/config"
	"fluxera/internal/handlers"
	"fluxera/internal/middleware"
	"fluxera/internal/repositories"
	"fluxera/internal/service"
	"fluxera/internal/storage"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := config.Load()

	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL is empty")
	}
	if cfg.JWTSecret == "" {
		log.Fatal("JWT_SECRET is empty")
	}

	store, err := storage.NewPostgresStore(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()

	userRepo := repositories.NewUserRepository(store.DB())

	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	authHandler := handlers.NewAuthHandler(authService)

	userService := service.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	r := chi.NewRouter()

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	r.Post("/auth/register", authHandler.Register)
	r.Post("/auth/login", authHandler.Login)

	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth(cfg.JWTSecret))
		r.Get("/me", userHandler.Me)
	})
	log.Printf("server started on %s", cfg.HTTPAddr)
	log.Fatal(http.ListenAndServe(cfg.HTTPAddr, r))
}
