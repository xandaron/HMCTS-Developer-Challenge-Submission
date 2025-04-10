package api

import (
	"HMCTS-Developer-Challenge/database"
	"HMCTS-Developer-Challenge/errors"
	"encoding/json"
	"net/http"
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

func HandleGetTasks(w http.ResponseWriter, r *http.Request, userID uint) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	dbHandle, err := db.GetDBHandle()
	if err != nil {
		errors.HandleServerError(w, err, "task.go: HandleGetTasks - GetDBHandle")
		return
	}

	rows, err := dbHandle.Query("SELECT * FROM tasks WHERE user_id = ?", userID)
	if err != nil {
		errors.HandleServerError(w, err, "task.go: HandleGetTasks - Query")
		return
	}

	var tasks []Task
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.UserID, &task.Name, &task.Description, &task.Status, &task.CreatedAt, &task.Deadline)
		if err != nil {
			errors.HandleServerError(w, err, "task.go: HandleGetTasks - Scan")
			return
		}
		tasks = append(tasks, task)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		errors.HandleServerError(w, err, "task.go: HandleTasks - Encode")
		return
	}
	w.WriteHeader(http.StatusOK)
}

func HandleAddTask(w http.ResponseWriter, r *http.Request, userID uint) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var jsonData struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Status      string `json:"status"`
		Deadline    string `json:"deadline"`
	}

	if err := json.NewDecoder(r.Body).Decode(&jsonData); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	dbHandle, err := db.GetDBHandle()
	if err != nil {
		errors.HandleServerError(w, err, "task.go: HandleAddTask - GetDBHandle")
		return
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
		errors.HandleServerError(w, err, "task.go: HandleAddTask - Exec")
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func HandleUpdateTask(w http.ResponseWriter, r *http.Request, userID uint) {

}

func HandleDeleteTask(w http.ResponseWriter, r *http.Request, userID uint) {

}
