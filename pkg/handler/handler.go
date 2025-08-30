package handler

import (
	"encoding/json"
	"net/http"
	"project_go-task-queue/pkg/models"
	"project_go-task-queue/pkg/worker"
)

type Handler struct {
	TaskQueue *worker.TaskQueue
}

func NewHandler(tq *worker.TaskQueue) *Handler {
	return &Handler{TaskQueue: tq}
}

func (h *Handler) Enqueue(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var t models.Task
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if t.ID == "" || t.Payload == "" {
		http.Error(w, "id and payload cannot be empty", http.StatusBadRequest)
		return
	}

	t.Status = "queued"
	h.TaskQueue.States.Store(t.ID, &t)

	select {
	case h.TaskQueue.Queue <- &t:
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Task queued"))
	default:
		http.Error(w, "Queue full", http.StatusServiceUnavailable)
	}
}

func (h *Handler) Tasks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
		return
	}

	tasks := []models.Task{}

	h.TaskQueue.States.Range(func(key, value any) bool {
		if t, ok := value.(*models.Task); ok {
			tasks = append(tasks, *t)
		}
		return true
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func (h *Handler) Healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
