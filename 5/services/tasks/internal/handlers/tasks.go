package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"

	"app/services/tasks/internal/repository"
	"app/services/tasks/internal/service"
)

type TaskHandler struct {
	repo repository.TaskRepository
	log  *logrus.Logger
}

func NewTaskHandler(
	repo repository.TaskRepository, log *logrus.Logger,
) *TaskHandler {
	return &TaskHandler{repo: repo, log: log}
}

func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		DueDate     string `json:"due_date"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid json"}`, http.StatusBadRequest)
		return
	}
	if req.Title == "" {
		http.Error(w, `{"error":"title required"}`, http.StatusBadRequest)
		return
	}

	task := service.Task{
		Title:       req.Title,
		Description: req.Description,
		DueDate:     req.DueDate,
		Done:        false,
	}

	task, err := h.repo.Create(task)
	if err != nil {
		h.log.Printf("failed to create task: %v", err)
		http.Error(
			w, `{"error":"internal error"}`, http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) List(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.repo.GetAll()
	if err != nil {
		h.log.Printf("failed to list tasks: %v", err)
		http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func (h *TaskHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/tasks/")
	if id == "" {
		http.Error(w, `{"error":"missing id"}`, http.StatusBadRequest)
		return
	}

	task, err := h.repo.GetByID(id)
	if err != nil {
		if err.Error() == "task not found" {
			http.Error(w, `{"error":"task not found"}`, http.StatusNotFound)
		} else {
			h.log.Printf("failed to get task: %v", err)
			http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/tasks/")
	if id == "" {
		http.Error(w, `{"error":"missing id"}`, http.StatusBadRequest)
		return
	}

	existing, err := h.repo.GetByID(id)
	if err != nil {
		if err.Error() == "task not found" {
			http.Error(w, `{"error":"task not found"}`, http.StatusNotFound)
		} else {
			h.log.Printf("failed to get task for update: %v", err)
			http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
		}
		return
	}

	var req struct {
		Title       *string `json:"title"`
		Description *string `json:"description"`
		DueDate     *string `json:"due_date"`
		Done        *bool   `json:"done"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid json"}`, http.StatusBadRequest)
		return
	}

	if req.Title != nil {
		existing.Title = *req.Title
	}
	if req.Description != nil {
		existing.Description = *req.Description
	}
	if req.DueDate != nil {
		existing.DueDate = *req.DueDate
	}
	if req.Done != nil {
		existing.Done = *req.Done
	}

	if err := h.repo.Update(existing); err != nil {
		if err.Error() == "task not found" {
			http.Error(w, `{"error":"task not found"}`, http.StatusNotFound)
		} else {
			h.log.Printf("failed to update task: %v", err)
			http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(existing)
}

func (h *TaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/tasks/")
	if id == "" {
		http.Error(w, `{"error":"missing id"}`, http.StatusBadRequest)
		return
	}

	if err := h.repo.Delete(id); err != nil {
		if err.Error() == "task not found" {
			http.Error(w, `{"error":"task not found"}`, http.StatusNotFound)
		} else {
			h.log.Printf("failed to delete task: %v", err)
			http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *TaskHandler) Search(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("title")
	if title == "" {
		http.Error(w, `{"error":"missing title"}`, http.StatusBadRequest)
		return
	}

	tasks, err := h.repo.SearchByTitle(title)
	if err != nil {
		h.log.Printf("search error: %v", err)
		http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func (h *TaskHandler) SearchVulnerable(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("title")
	if title == "" {
		http.Error(w, `{"error":"missing title"}`, http.StatusBadRequest)
		return
	}

	tasks, err := h.repo.SearchByTitleVulnerable(title)
	if err != nil {
		h.log.Printf("search error: %v", err)
		http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}
