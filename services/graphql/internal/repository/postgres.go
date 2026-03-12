package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"

	"app/services/graphql/internal/models"
)

type PostgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) *PostgresRepo {
	return &PostgresRepo{db: db}
}

func (r *PostgresRepo) Create(task models.Task) error {
	task.ID = uuid.New().String()
	task.CreatedAt = time.Now()
	task.UpdatedAt = task.CreatedAt

	var dueDate sql.NullString
	if task.DueDate != nil && *task.DueDate != "" {
		dueDate.String = *task.DueDate
		dueDate.Valid = true
	}

	query := `INSERT INTO tasks (id, title, description, done, due_date, created_at, updated_at)
              VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.Exec(query, task.ID, task.Title, task.Description, task.Done, dueDate, task.CreatedAt, task.UpdatedAt)
	return err
}

func (r *PostgresRepo) GetAll() ([]models.Task, error) {
	rows, err := r.db.Query(`SELECT id, title, description, done, due_date, created_at, updated_at FROM tasks`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var t models.Task
		var dueDate sql.NullString
		err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Done, &dueDate, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			return nil, err
		}
		if dueDate.Valid {
			t.DueDate = &dueDate.String
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (r *PostgresRepo) GetByID(id string) (models.Task, error) {
	var t models.Task
	var dueDate sql.NullString
	query := `SELECT id, title, description, done, due_date, created_at, updated_at FROM tasks WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(&t.ID, &t.Title, &t.Description, &t.Done, &dueDate, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return t, fmt.Errorf("task not found")
		}
		return t, err
	}
	if dueDate.Valid {
		t.DueDate = &dueDate.String
	}
	return t, nil
}

func (r *PostgresRepo) Update(task models.Task) error {
	task.UpdatedAt = time.Now()

	var dueDate sql.NullString
	if task.DueDate != nil && *task.DueDate != "" {
		dueDate.String = *task.DueDate
		dueDate.Valid = true
	}

	query := `UPDATE tasks SET title=$1, description=$2, done=$3, due_date=$4, updated_at=$5 WHERE id=$6`
	result, err := r.db.Exec(query, task.Title, task.Description, task.Done, dueDate, task.UpdatedAt, task.ID)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("task not found")
	}
	return nil
}

func (r *PostgresRepo) Delete(id string) error {
	query := `DELETE FROM tasks WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("task not found")
	}
	return nil
}
