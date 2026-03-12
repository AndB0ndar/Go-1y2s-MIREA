package models

import "time"

type Task struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description *string   `json:"description"`
	Done        bool      `json:"done"`
	DueDate     *string   `json:"dueDate"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
