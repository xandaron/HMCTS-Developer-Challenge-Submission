package api

import (
	"HMCTS-Developer-Challenge/database"
	"HMCTS-Developer-Challenge/errors"
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
)

type Task struct {
	ID          uint   `json:"id"`
	UserID      uint   `json:"user_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
	Deadline    string `json:"deadline"`
}

func HandleTasks(w http.ResponseWriter, r *http.Request, userID uint) {
	switch r.Method {
	case http.MethodGet:
		tasks, err := getTasks(userID)
		if err != nil {
			errors.HandleServerError(w, err, "task.go: HandleTasks - getTasks")
			break
		}

		var buf bytes.Buffer
		if err := json.NewEncoder(&buf).Encode(tasks); err != nil {
			errors.HandleServerError(w, err, "task.go: HandleTasks - Encode")
			break
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		buf.WriteTo(w)
	case http.MethodPost:
		err := addTask(r, userID)
		if err != nil {
			errors.HandleServerError(w, err, "task.go: HandleTasks - addTask")
			break
		}

		w.WriteHeader(http.StatusCreated)
	case http.MethodPut:
		editTask(w, r, userID)
	case http.MethodDelete:
		pathParts := strings.Split(r.URL.Path, "/")
		if len(pathParts) < 4 || pathParts[3] == "" {
			http.Error(w, "Task ID is required", http.StatusBadRequest)
		}

		if err := deleteTask(userID, pathParts[3]); err != nil {
			errors.HandleServerError(w, err, "task.go: HandleTasks - deleteTask")
			break
		}

		w.WriteHeader(http.StatusNoContent)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getTasks(userID uint) ([]Task, error) {
	dbHandle, err := db.GetDBHandle()
	if err != nil {
		return nil, errors.New(err, "task.go: HandleGetTasks - GetDBHandle")
	}

	rows, err := dbHandle.Query("SELECT * FROM tasks WHERE user_id = ?", userID)
	if err != nil {
		return nil, errors.New(err, "task.go: HandleGetTasks - Query")
	}

	var tasks []Task
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.UserID, &task.Name, &task.Description, &task.Status, &task.CreatedAt, &task.Deadline)
		if err != nil {
			return nil, errors.New(err, "task.go: HandleGetTasks - Scan")
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func addTask(r *http.Request, userID uint) error {
	var jsonData struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Status      string `json:"status"`
		Deadline    string `json:"deadline"`
	}
	if err := json.NewDecoder(r.Body).Decode(&jsonData); err != nil {
		return errors.New(err, "task.go: addTask - Decode")
	}

	dbHandle, err := db.GetDBHandle()
	if err != nil {
		return errors.New(err, "task.go: addTask - GetDBHandle")
	}

	_, err = dbHandle.Exec(
		"INSERT INTO tasks (user_id, name, description, status, deadline) VALUES (?, ?, ?, ?, ?)",
		userID,
		jsonData.Name,
		jsonData.Description,
		jsonData.Status,
		jsonData.Deadline,
	)
	if err != nil {
		return errors.New(err, "task.go: addTask - Exec")
	}

	return nil
}

func editTask(w http.ResponseWriter, r *http.Request, userID uint) {

}

func deleteTask(userID uint, taskID string) error {
	dbHandle, err := db.GetDBHandle()
	if err != nil {
		return errors.New(err, "task.go: deleteTask - GetDBHandle")
	}

	_, err = dbHandle.Exec("DELETE FROM tasks WHERE id = ? AND user_id = ?", taskID, userID)
	if err != nil {
		return errors.New(err, "task.go: deleteTask - Exec")
	}

	return nil
}
