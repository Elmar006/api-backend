package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Elmar006/api-backend/internal/auth"
	"github.com/Elmar006/api-backend/internal/database"
	"github.com/Elmar006/api-backend/internal/models"
	"github.com/go-chi/chi/v5"
)

func GetTasks(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.GetUserId(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
	}

	log.Printf("GetTasks called by user_id: %d", userID)

	tasks := []models.Task{}
	resultDB := database.DB.Db.Where("user_id = ?", userID).Find(&tasks)
	if resultDB.Error != nil {
		http.Error(w, resultDB.Error.Error(), http.StatusBadRequest)
		return
	}

	resp, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)

	log.Printf("GetTasks: returned %d tasks for user_id %d", len(tasks), userID)
}

func GetTaskByID(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.GetUserId(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		// Если ID не число
		http.Error(w,
			fmt.Sprintf("Invalid task ID: %s", idParam),
			http.StatusBadRequest)
		return
	}

	if id <= 0 {
		http.Error(w, "Task ID must be positive", http.StatusBadRequest)
		return
	}

	task := new(models.Task)
	resultDB := database.DB.Db.
		Where("id = ? AND user_id = ?", id, userID).
		First(&task)

	if resultDB.Error != nil {
		if resultDB.Error.Error() == "record not found" {
			http.Error(w,
				fmt.Sprintf("Task with ID %d not found", id),
				http.StatusNotFound)
		} else {
			http.Error(w,
				fmt.Sprintf("Database error: %v", resultDB.Error),
				http.StatusBadRequest)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
	log.Printf("GetTaskByID: returned task_id=%d for user_id=%d", id, userID)
}

func PostTask(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.GetUserId(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	task := new(models.Task)
	var buf bytes.Buffer
	_, err = buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(buf.Bytes(), &task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	task.UserID = userID
	resultDB := database.DB.Db.Create(&task)
	if resultDB.Error != nil {
		http.Error(w, resultDB.Error.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Location", fmt.Sprintf("/task/%d", task.ID))
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
	log.Printf("PostTask: created task_id=%d for user_id=%d", task.ID, userID)
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.GetUserId(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w,
			fmt.Sprintf("Invalid task ID: %v", err),
			http.StatusBadRequest)
		return
	}

	resultDB := database.DB.Db.
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&models.Task{})

	if resultDB.Error != nil {
		http.Error(w,
			fmt.Sprintf("Database error: %v", resultDB.Error),
			http.StatusBadRequest)
		return
	}

	if resultDB.RowsAffected == 0 {
		http.Error(w,
			fmt.Sprintf("Task with ID %d not found or access denied", id),
			http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Task deleted successfully",
		"task_id": id,
	})

	log.Printf("DeleteTask: deleted task_id=%d by user_id=%d", id, userID)
}
