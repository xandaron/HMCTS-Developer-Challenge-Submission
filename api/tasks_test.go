package api

import (
	"HMCTS-Developer-Challenge/database"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	err := database.Connect()
	if err != nil {
		fmt.Printf("Failed to connect to database: %v\n", err)
		os.Exit(1)
	}
	exitCode := m.Run()
	database.Disconnect()

	os.Exit(exitCode)
}

func TestGetTasks(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/tasks", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		TasksHandler(w, r, 1)
	})

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check that the response is a JSON array
	var tasks []task
	err = json.Unmarshal(rr.Body.Bytes(), &tasks)
	if err != nil {
		t.Errorf("Failed to parse response as JSON: %v", err)
	}

	// Check that at least we have the tasks from the seed data for user 1
	if len(tasks) < 2 {
		t.Errorf("Expected at least 2 tasks for user 1, got %d", len(tasks))
	}
}

func TestGetSingleTask(t *testing.T) {
	// Get task with ID 1 for user 1
	req, err := http.NewRequest("GET", "/api/tasks/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		TasksHandler(w, r, 1) // User ID 1
	})

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check that the response is a JSON object
	var taskData task
	err = json.Unmarshal(rr.Body.Bytes(), &taskData)
	if err != nil {
		t.Errorf("Failed to parse response as JSON: %v", err)
	}

	// Check the task ID and user ID
	if taskData.ID != 1 {
		t.Errorf("Expected task ID 1, got %d", taskData.ID)
	}
	if taskData.UserID != 1 {
		t.Errorf("Expected user ID 1, got %d", taskData.UserID)
	}
}

func TestGetTaskNotFound(t *testing.T) {
	// Try to get a non-existent task
	req, err := http.NewRequest("GET", "/api/tasks/999", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		TasksHandler(w, r, 1) // User ID 1
	})

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}

func TestAddTask(t *testing.T) {
	data := jsonData{
		Name:        "Test Task",
		Description: "This is a test task",
		Status:      "INCOMPLETE",
		Deadline:    "2025-12-31 00:00:00",
	}
	taskJSON, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/api/tasks", bytes.NewBuffer(taskJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		TasksHandler(w, r, 1)
	})

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	req, err = http.NewRequest("GET", "/api/tasks", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	var tasks []task
	err = json.Unmarshal(rr.Body.Bytes(), &tasks)
	if err != nil {
		t.Errorf("Failed to parse response as JSON: %v", err)
	}

	found := false
	for _, t := range tasks {
		if t.Name == "Test Task" && t.Description == "This is a test task" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Added task was not found in the list of tasks")
	}
}

func TestAddTaskMissingData(t *testing.T) {
	task := jsonData{
		Name:        "",
		Description: "This is a test task",
		Status:      "INCOMPLETE",
		Deadline:    "2025-12-31 00:00:00",
	}
	taskJSON, err := json.Marshal(task)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/api/tasks", bytes.NewBuffer(taskJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		TasksHandler(w, r, 1)
	})

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestEditTask(t *testing.T) {
	// Updated task data
	editedTask := jsonData{
		Name:        "Edited Task",
		Description: "This task has been edited",
		Status:      "COMPLETE",
		Deadline:    "2025-11-30 00:00:00",
	}
	taskJSON, err := json.Marshal(editedTask)
	if err != nil {
		t.Fatal(err)
	}

	// Edit the task
	req, err := http.NewRequest("PUT", "/api/tasks/1", bytes.NewBuffer(taskJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		TasksHandler(w, r, 1) // User ID 1
	})

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
	}

	// Now check if the task was actually updated
	req, err = http.NewRequest("GET", "/api/tasks/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	var taskData task
	err = json.Unmarshal(rr.Body.Bytes(), &taskData)
	if err != nil {
		t.Errorf("Failed to parse response as JSON: %v", err)
	}

	// Check if the task was updated
	if taskData.Name != "Edited Task" {
		t.Errorf("Expected task name 'Edited Task', got '%s'", taskData.Name)
	}
	if taskData.Description != "This task has been edited" {
		t.Errorf("Expected task description 'This task has been edited', got '%s'", taskData.Description)
	}
	if taskData.Status != "COMPLETE" {
		t.Errorf("Expected task status 'COMPLETE', got '%s'", taskData.Status)
	}
}

func TestEditTaskNotFound(t *testing.T) {
	// Try to edit a non-existent task
	taskData := jsonData{
		Name:        "Edited Task",
		Description: "This task has been edited",
		Status:      "COMPLETE",
		Deadline:    "2025-11-30 00:00:00",
	}
	taskJSON, err := json.Marshal(taskData)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("PUT", "/api/tasks/999", bytes.NewBuffer(taskJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		TasksHandler(w, r, 1) // User ID 1
	})

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestDeleteTask(t *testing.T) {
	// First, add a task to delete
	addTask := jsonData{
		Name:        "Task to Delete",
		Description: "This task will be deleted",
		Status:      "INCOMPLETE",
		Deadline:    "2025-12-31 00:00:00",
	}
	taskJSON, err := json.Marshal(addTask)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/api/tasks", bytes.NewBuffer(taskJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		TasksHandler(w, r, 1) // User ID 1
	})

	handler.ServeHTTP(rr, req)

	// Now get all tasks to find the ID of the task we just added
	req, err = http.NewRequest("GET", "/api/tasks", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	var tasks []task
	err = json.Unmarshal(rr.Body.Bytes(), &tasks)
	if err != nil {
		t.Errorf("Failed to parse response as JSON: %v", err)
	}

	// Find the task we just added
	var taskIDToDelete string
	for _, task := range tasks {
		if task.Name == "Task to Delete" {
			taskIDToDelete = fmt.Sprintf("%d", task.ID)
			break
		}
	}

	if taskIDToDelete == "" {
		t.Fatal("Could not find the task to delete")
	}

	// Now delete the task
	req, err = http.NewRequest("DELETE", "/api/tasks/"+taskIDToDelete, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
	}

	// Verify the task is gone
	req, err = http.NewRequest("GET", "/api/tasks/"+taskIDToDelete, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}

func TestDeleteTaskNotFound(t *testing.T) {
	// Try to delete a non-existent task
	req, err := http.NewRequest("DELETE", "/api/tasks/999", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		TasksHandler(w, r, 1) // User ID 1
	})

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}

func TestTasksHandlerHttpMethods(t *testing.T) {
	// Test OPTIONS method
	req, err := http.NewRequest("OPTIONS", "/api/tasks", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		TasksHandler(w, r, 1) // User ID 1
	})

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code for OPTIONS: got %v want %v", status, http.StatusNoContent)
	}

	// Test unsupported method
	req, err = http.NewRequest("PATCH", "/api/tasks", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code for unsupported method: got %v want %v", status, http.StatusMethodNotAllowed)
	}
}

func TestGetTaskBelongingToAnotherUser(t *testing.T) {
	// Try to get task 3 as user 1 (task 3 belongs to user 2)
	req, err := http.NewRequest("GET", "/api/tasks/3", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		TasksHandler(w, r, 1) // User ID 1 trying to access user 2's task
	})

	handler.ServeHTTP(rr, req)

	// Should return not found, as the task exists but doesn't belong to user 1
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}

func TestEditTaskBelongingToAnotherUser(t *testing.T) {
	// Try to edit task 3 as user 1 (task 3 belongs to user 2)
	taskData := jsonData{
		Name:        "Edited Task That Belongs To Another User",
		Description: "This task belongs to user 2",
		Status:      "COMPLETE",
		Deadline:    "2025-11-30 00:00:00",
	}
	taskJSON, err := json.Marshal(taskData)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("PUT", "/api/tasks/3", bytes.NewBuffer(taskJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		TasksHandler(w, r, 1) // User ID 1 trying to edit user 2's task
	})

	handler.ServeHTTP(rr, req)

	// Should return bad request with task not found error
	if rr.Code != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, http.StatusBadRequest)
	}

	// Verify the task was not modified by checking it as user 2
	req, err = http.NewRequest("GET", "/api/tasks/3", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		TasksHandler(w, r, 2) // User ID 2 (owner of the task)
	})

	handler.ServeHTTP(rr, req)

	var taskInfo task
	err = json.Unmarshal(rr.Body.Bytes(), &taskInfo)
	if err != nil {
		t.Errorf("Failed to parse response as JSON: %v", err)
	}

	// Task should still have its original name
	if taskInfo.Name != "Task 3" {
		t.Errorf("Task was modified when it shouldn't have been. Expected name 'Task 3', got '%s'", taskInfo.Name)
	}
}

func TestDeleteTaskBelongingToAnotherUser(t *testing.T) {
	// Try to delete task 3 as user 1 (task 3 belongs to user 2)
	req, err := http.NewRequest("DELETE", "/api/tasks/3", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		TasksHandler(w, r, 1) // User ID 1 trying to delete user 2's task
	})

	handler.ServeHTTP(rr, req)

	// Should return not found since the task exists but doesn't belong to user 1
	if rr.Code != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, http.StatusNotFound)
	}

	// Verify the task still exists by checking it as user 2
	req, err = http.NewRequest("GET", "/api/tasks/3", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		TasksHandler(w, r, 2) // User ID 2 (owner of the task)
	})

	handler.ServeHTTP(rr, req)

	// Should be able to find the task
	if rr.Code != http.StatusOK {
		t.Errorf("Task was deleted when it shouldn't have been. Expected status OK, got %v", rr.Code)
	}
}
