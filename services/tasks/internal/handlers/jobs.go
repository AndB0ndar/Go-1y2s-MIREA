package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

type JobHandler struct {
	rabbitConn *amqp.Connection
	queueName  string
	log        *logrus.Logger
}

func NewJobHandler(conn *amqp.Connection, queue string, log *logrus.Logger) *JobHandler {
	return &JobHandler{
		rabbitConn: conn,
		queueName:  queue,
		log:        log,
	}
}

type ProcessTaskRequest struct {
	TaskID string `json:"task_id"`
}

func (h *JobHandler) ProcessTask(w http.ResponseWriter, r *http.Request) {
	var req ProcessTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid json"}`, http.StatusBadRequest)
		return
	}
	if req.TaskID == "" {
		http.Error(w, `{"error":"task_id required"}`, http.StatusBadRequest)
		return
	}

	messageID := uuid.New().String()
	job := map[string]interface{}{
		"job":        "process_task",
		"task_id":    req.TaskID,
		"attempt":    1,
		"message_id": messageID,
		"timestamp":  time.Now().UTC().Format(time.RFC3339),
	}
	body, err := json.Marshal(job)
	if err != nil {
		h.log.WithError(err).Error("failed to marshal job")
		http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
		return
	}

	ch, err := h.rabbitConn.Channel()
	if err != nil {
		h.log.WithError(err).Error("failed to open channel")
		http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
		return
	}
	defer ch.Close()

	err = ch.Publish(
		"",          // exchange
		h.queueName, // routing key (queue name)
		true,        // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
			MessageId:    messageID,
		})
	if err != nil {
		h.log.WithError(err).Error("failed to publish job")
		http.Error(w, `{"error":"internal error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{
		"status":     "accepted",
		"message_id": messageID,
	})
}
