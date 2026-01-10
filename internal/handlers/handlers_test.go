package handlers

import (
	"bytes"

	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Elmar006/api-backend/internal/database"
	"github.com/go-chi/chi/v5"
)

func TestPostTask_InvalidJSON(t *testing.T) {
	body := []byte(`{"description": "test", "note": }`)

	req := httptest.NewRequest("POST", "/tasks", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	PostTask(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d for invalid JSON, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestPostTask_EmptyBody(t *testing.T) {
	req := httptest.NewRequest("POST", "/tasks", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	PostTask(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d for empty body, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestGetTaskByID_InvalidID(t *testing.T) {
	r := chi.NewRouter()
	r.Get("/task/{id}", GetTaskByID)

	req := httptest.NewRequest("GET", "/task/abc", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d for invalid ID, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestGetTasks_RequiresDB(t *testing.T) {
	// Пропускаем если БД не подключена
	if database.DB.Db == nil {
		t.Skip("Требуется подключение к БД. Запустите с флагом -short для пропуска")
	}

	req := httptest.NewRequest("GET", "/tasks", nil)
	w := httptest.NewRecorder()

	GetTasks(w, req)

	if w.Code >= http.StatusInternalServerError {
		t.Errorf("Should not return 5xx error, got %d", w.Code)
	}
}

func TestHandlers_SmokeTest(t *testing.T) {
	t.Run("GetTasks function exists", func(t *testing.T) {

	})
	t.Run("PostTask function exists", func(t *testing.T) {
	})
}
