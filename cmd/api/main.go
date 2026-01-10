package main

import (
	"fmt"
	"net/http"

	"github.com/Elmar006/api-backend/internal/database"
	"github.com/Elmar006/api-backend/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func main() {
	database.ConnectDB()
	r := chi.NewRouter()

	r.Get("/tasks", handlers.GetTasks)
	r.Post("/tasks", handlers.PostTask)
	r.Get("/task/{id}", handlers.GetTaskByID)
	r.Delete("/task/{id}", handlers.DeleteTask)

	if err := http.ListenAndServe(":3000", r); err != nil {
		fmt.Printf("Error start serve: %v", err)
		return
	}
}
