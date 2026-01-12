package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Elmar006/api-backend/internal/auth"
	"github.com/Elmar006/api-backend/internal/database"
	"github.com/Elmar006/api-backend/internal/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	database.ConnectDB()
	auth.InitJWT()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	r.Post("/register", handlers.Register)
	r.Post("/login", handlers.Login)

	r.Group(func(r chi.Router) {
		r.Use(auth.Verifier())
		r.Use(auth.Authenticator())

		r.Get("/tasks", handlers.GetTasks)
		r.Post("/tasks", handlers.PostTask)
		r.Get("/task/{id}", handlers.GetTaskByID)
		r.Delete("/task/{id}", handlers.DeleteTask)
		r.Get("/me", handlers.GetCurrentUser)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Server starting on http://0.0.0.0:%s", port)

	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
