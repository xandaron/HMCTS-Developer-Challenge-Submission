package api

import (
	"HMCTS-Developer-Challenge/database"
	"HMCTS-Developer-Challenge/errors"
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
)

type task struct {
	ID          uint   `json:"id"`
	UserID      uint   `json:"user_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
	Deadline    string `json:"deadline"`
}

type jsonData struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Deadline    string `json:"deadline"`
}

var errTaskNotFound = errors.Error("Task Not Found")
var errMissingJsonData = errors.Error("Missing JSON Data")

func TasksHandler(w http.ResponseWriter, r *http.Request, userID uint) {
	switch r.Method {
	case http.MethodGet:
		var tasks any
		var err error

		pathParts := strings.Split(r.URL.Path, "/")
		if len(pathParts) < 4 || pathParts[3] == "" {
			tasks, err = getTasks(userID)
		} else {
			tasks, err = getTask(userID, pathParts[3])
			if err == errTaskNotFound {
				http.Error(w, errTaskNotFound.Error(), http.StatusNotFound)
				break
			}
		}

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
		var data jsonData
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			break
		}

		if err := addTask(userID, data); err == errMissingJsonData {
			http.Error(w, "Missing JSON data", http.StatusBadRequest)
			break
		} else if err != nil {
			errors.HandleServerError(w, err, "task.go: HandleTasks - addTask")
			break
		}

		w.WriteHeader(http.StatusCreated)
	case http.MethodPut:
		pathParts := strings.Split(r.URL.Path, "/")
		if len(pathParts) < 4 || pathParts[3] == "" {
			http.Error(w, "Task ID Required", http.StatusBadRequest)
			break
		}

		var data jsonData
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			break
		}

		if err := editTask(userID, pathParts[3], data); err == errMissingJsonData || err == errTaskNotFound {
			http.Error(w, err.Error(), http.StatusBadRequest)
			break
		} else if err != nil {
			errors.HandleServerError(w, err, "task.go: HandleTasks - editTask")
			break
		}

		w.WriteHeader(http.StatusNoContent)
	case http.MethodDelete:
		pathParts := strings.Split(r.URL.Path, "/")
		if len(pathParts) < 4 || pathParts[3] == "" {
			http.Error(w, "Task ID Required", http.StatusBadRequest)
			break
		}

		if err := deleteTask(userID, pathParts[3]); err == errTaskNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			break
		} else if err != nil {
			errors.HandleServerError(w, err, "task.go: HandleTasks - deleteTask")
			break
		}

		w.WriteHeader(http.StatusNoContent)
	case http.MethodOptions:
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.WriteHeader(http.StatusNoContent)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getTask(userID uint, taskID string) (task, error) {
	var t task

	dbHandle, err := db.GetDBHandle()
	if err != nil {
		return t, errors.AddContext(err, "task.go: HandleGetTask - GetDBHandle")
	}

	if err := dbHandle.QueryRow(
		"SELECT * FROM tasks WHERE id = ? AND user_id = ?", taskID, userID).Scan(
		&t.ID,
		&t.UserID,
		&t.Name,
		&t.Description,
		&t.Status,
		&t.CreatedAt,
		&t.Deadline,
	); err != nil {
		if err == sql.ErrNoRows {
			return t, errTaskNotFound
		}
		return t, errors.AddContext(err, "task.go: HandleGetTask - QueryRow")
	}
	return t, nil
}

func getTasks(userID uint) ([]task, error) {
	dbHandle, err := db.GetDBHandle()
	if err != nil {
		return nil, errors.AddContext(err, "task.go: HandleGetTasks - GetDBHandle")
	}

	rows, err := dbHandle.Query("SELECT * FROM tasks WHERE user_id = ?", userID)
	if err != nil {
		return nil, errors.AddContext(err, "task.go: HandleGetTasks - Query")
	}
	defer rows.Close()

	var tasks []task
	for rows.Next() {
		var t task
		err := rows.Scan(&t.ID, &t.UserID, &t.Name, &t.Description, &t.Status, &t.CreatedAt, &t.Deadline)
		if err != nil {
			return nil, errors.AddContext(err, "task.go: HandleGetTasks - Scan")
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func addTask(userID uint, data jsonData) error {
	dbHandle, err := db.GetDBHandle()
	if err != nil {
		return errors.AddContext(err, "task.go: addTask - GetDBHandle")
	}

	if data.Name == "" || data.Status == "" || data.Deadline == "" {
		return errMissingJsonData
	}

	_, err = dbHandle.Exec(
		"INSERT INTO tasks (user_id, name, description, status, deadline) VALUES (?, ?, ?, ?, ?)",
		userID,
		data.Name,
		data.Description,
		data.Status,
		data.Deadline,
	)
	if err != nil {
		return errors.AddContext(err, "task.go: addTask - Exec")
	}

	return nil
}

func editTask(userID uint, taskID string, data jsonData) error {
	dbHandle, err := db.GetDBHandle()
	if err != nil {
		return errors.AddContext(err, "task.go: editTask - GetDBHandle")
	}

	if exists, err := checkTaskExists(userID, taskID); err != nil {
		return errors.AddContext(err, "task.go: deleteTask - QueryRow")
	} else if !exists {
		return errTaskNotFound
	}

	if data.Name == "" || data.Status == "" || data.Deadline == "" {
		return errMissingJsonData
	}

	_, err = dbHandle.Exec(
		"UPDATE tasks SET name = ?, description = ?, status = ?, deadline = ? WHERE id = ? AND user_id = ?",
		data.Name,
		data.Description,
		data.Status,
		data.Deadline,
		taskID,
		userID,
	)
	if err != nil {
		return errors.AddContext(err, "task.go: editTask - Exec")
	}

	return nil
}

func deleteTask(userID uint, taskID string) error {
	dbHandle, err := db.GetDBHandle()
	if err != nil {
		return errors.AddContext(err, "task.go: deleteTask - GetDBHandle")
	}

	if exists, err := checkTaskExists(userID, taskID); err != nil {
		return errors.AddContext(err, "task.go: deleteTask - QueryRow")
	} else if !exists {
		return errTaskNotFound
	}

	_, err = dbHandle.Exec("DELETE FROM tasks WHERE id = ? AND user_id = ?", taskID, userID)
	if err != nil {
		return errors.AddContext(err, "task.go: deleteTask - Exec")
	}

	return nil
}

func checkTaskExists(userID uint, taskID string) (bool, error) {
	dbHandle, err := db.GetDBHandle()
	if err != nil {
		return false, errors.AddContext(err, "task.go: checkTaskExists - GetDBHandle")
	}

	var exists bool
	if err := dbHandle.QueryRow("SELECT EXISTS(SELECT 1 FROM tasks WHERE id = ? AND user_id = ?)", taskID, userID).Scan(&exists); err != nil {
		return false, errors.AddContext(err, "task.go: checkTaskExists - QueryRow")
	}
	return exists, nil
}
