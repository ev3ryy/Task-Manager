package tasks

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"restful/src/database"
	"restful/src/types"
)

func CreateTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "this method is not supported!", http.StatusMethodNotAllowed)
		return
	}

	var task types.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		http.Error(w, "required fields are missing: title", http.StatusBadRequest)
		return
	}

	insertQuery := `
		INSERT INTO Tasks (title, description)
		VALUES ($1, $2)
		RETURNING id;
	`
	var id int
	err := database.DB.QueryRow(insertQuery, task.Title, task.Description).Scan(&id)
	if err != nil {
		http.Error(w, "failed to create task", http.StatusInternalServerError)
		return
	}

	res := types.CreateTaskResponse{
		Title:   task.Title,
		Message: "Задача успешно создана!",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func GetTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "this method is not supported!", http.StatusMethodNotAllowed)
		return
	}

	title := r.URL.Query().Get("title")
	if title == "" {
		http.Error(w, "missing title parameter", http.StatusBadRequest)
		return
	}

	selectQuery := `
		SELECT id, title, description, completed
		FROM Tasks
		WHERE title = $1;
	`

	var task types.Task

	err := database.DB.QueryRow(selectQuery, title).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Completed,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "task not found", http.StatusNotFound)
		}
		return
	}

	res := types.GetTaskResponse{
		Message: task.Description,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func GetAllTasks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "this method is not supported!", http.StatusMethodNotAllowed)
		return
	}

	selectQuery := `
		SELECT id, title, description, completed
		FROM Tasks
		ORDER BY id;
	`

	rows, err := database.DB.Query(selectQuery)
	if err != nil {
		http.Error(w, "failed to retrieve tasks from database", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	tasks := []types.Task{}

	for rows.Next() {
		var task types.Task

		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Completed,
		)
		if err != nil {
			http.Error(w, "failed to proccess task data", http.StatusInternalServerError)
			return
		}

		tasks = append(tasks, task)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "this method is not supported!", http.StatusMethodNotAllowed)
		return
	}

	var task types.UpdateTask
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		http.Error(w, "title is required, but empty", http.StatusBadRequest)
		return
	}

	updateQuery := `
		UPDATE Tasks
		SET description = $1, completed = $2
		WHERE title = $3
		RETURNING id;
	`

	var updatedID int

	err := database.DB.QueryRow(updateQuery,
		task.Description,
		task.Completed,
		task.Title,
	).Scan(&updatedID)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "task not found", http.StatusNotFound)
		}
		return
	}

	response := types.UpdateTaskResponse{
		Title:       task.Title,
		Description: task.Description,
		Completed:   task.Completed,
		Message:     fmt.Sprintf("task '%s' (ID: %d) updated successfully.", task.Title, updatedID),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(response)
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "this method is not supported!", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "id parametr is required, but empty", http.StatusBadRequest)
		return
	}

	taskID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id format", http.StatusBadRequest)
		return
	}

	deleteQuery := `
		DELETE FROM Tasks
		WHERE id = $1
		RETURNING id;
	`

	var deletedTaskID int
	err = database.DB.QueryRow(deleteQuery, taskID).Scan(&deletedTaskID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Task not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to delete task from database", http.StatusInternalServerError)
		}
		return
	}

	res := types.DeleteTaskResponse{
		Message: fmt.Sprintf("Task with '%d' id was deleted", deletedTaskID),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(res)
}
