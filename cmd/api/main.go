package main

import (
	"context"
	"fluxera/internal/config"
	"fluxera/internal/events"
	"fluxera/internal/handlers"
	"fluxera/internal/middleware"
	"fluxera/internal/repositories"
	"fluxera/internal/service"
	"fluxera/internal/storage"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
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

	dispatcher := events.NewDispatcher(100)

	userRepo := repositories.NewUserRepository(store.DB())

	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	authHandler := handlers.NewAuthHandler(authService)

	userService := service.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	projectRepo := repositories.NewProjectRepository(store.DB())
	projectService := service.NewProjectService(projectRepo)
	projectHandler := handlers.NewProjectHandler(projectService)

	activityRepo := repositories.NewActivityLogRepository(store.DB())
	activityService := service.NewActivityLogService(activityRepo, projectRepo)
	eventActivityHandler := events.NewActivityHandler(activityService)
	dispatcher.Start(context.Background(), eventActivityHandler)
	activityHandler := handlers.NewActivityLogHandler(activityService)

	taskRepo := repositories.NewTaskRepository(store.DB())
	taskService := service.NewTaskService(taskRepo, projectRepo, dispatcher)
	taskHandler := handlers.NewTaskHandler(taskService)

	commentRepo := repositories.NewCommentRepository(store.DB())
	commentService := service.NewCommentService(commentRepo, taskRepo, projectRepo, dispatcher)
	commentHandler := handlers.NewCommentHandler(commentService)

	r := chi.NewRouter()
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

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
		r.Post("/projects", projectHandler.CreateProject)
		r.Get("/projects", projectHandler.GetProjects)
		r.Get("/projects/{id}", projectHandler.GetProjectByID)
		r.Delete("/projects/{id}", projectHandler.DeleteProject)
		r.Get("/projects/{projectID}/tasks", taskHandler.GetTasksByProject)
		r.Post("/projects/{projectID}/tasks", taskHandler.CreateTask)
		r.Put("/tasks/{id}", taskHandler.UpdateTask)
		r.Patch("/tasks/{id}/status", taskHandler.UpdateTaskStatus)
		r.Delete("/tasks/{id}", taskHandler.DeleteTask)
		r.Post("/tasks/{id}/comments", commentHandler.CreateComment)
		r.Get("/tasks/{id}/comments", commentHandler.GetCommentsByTask)
		r.Put("/comments/{id}", commentHandler.UpdateComment)
		r.Delete("/comments/{id}", commentHandler.DeleteComment)
		r.Get("/projects/{id}/activity", activityHandler.GetProjectActivity)
	})
	log.Printf("server started on %s", cfg.HTTPAddr)
	log.Fatal(http.ListenAndServe(cfg.HTTPAddr, r))
}
