package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/lahnasti/GO_praktikum/internal/domain/models"
)

type DBstorage struct {
	conn *pgx.Conn
}

func NewDB(conn *pgx.Conn) DBstorage {
	return DBstorage{
		conn: conn,
	}
}

func (db *DBstorage) GetAllTasks() ([]models.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rows, err := db.conn.Query(ctx, "SELECT id, title, description FROM tasks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Description); err != nil {
			return nil, err
		}
		task.Title = strings.TrimSpace(task.Title)
		task.Description = strings.TrimSpace(task.Description)
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (db *DBstorage) GetTasksByID(id string) (models.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	row := db.conn.QueryRow(ctx, "SELECT title, description FROM tasks WHERE id=$1", id)
	var task models.Task
	if err := row.Scan(&task.Title, &task.Description); err != nil {
		return models.Task{}, err
	}
	return task, nil
}

func (db *DBstorage) AddTask(task models.Task) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := "INSERT INTO users (title, descrtiption) VALUES ($1, $2) RETURNING id"
	var taskID string
	err := db.conn.QueryRow(ctx, query, task.Title, task.Description).Scan(&taskID)
	if err != nil {
		return "", fmt.Errorf("failed to insert task: %w", err)
	}
	// Проверка, что taskID не пустой
	if taskID == "" {
		return "", fmt.Errorf("failed to get userID after insert")
	}
	return taskID, nil
}

func (db *DBstorage) UpdateTask(id string, task models.Task) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := db.conn.Exec(ctx, "UPDATE tasks SET title=$1, description=$2 WHERE id=$3", task.Title, task.Description, id)
	if err != nil {
		return fmt.Errorf("update task failed: %w", err)
	}
	return nil
}

func (db *DBstorage) DeleteTask(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := db.conn.Exec(ctx, "DELETE FROM tasks WHERE id=$1", id)
	if err != nil {
		return fmt.Errorf("delete task failed: %w", err)
	}
	return nil
}