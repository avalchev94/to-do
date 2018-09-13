package main

import (
	"encoding/json"
	"errors"
	"net/http"
)

func handleTask(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		handleTaskPost(w, r)
	case "GET":
		handleTaskGet(w, r)
	case "DELETE":
		handleTaskDelete(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func handleTaskPost(w http.ResponseWriter, r *http.Request) {
	var task Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		respondErr(w, r, err, http.StatusBadRequest)
		return
	}

	// Get current logged user
	userID, err := getCurrentUserID(r)
	if err != nil {
		respondErr(w, r, errors.New("no logged user"), http.StatusBadRequest)
		return
	}
	task.UserID = userID

	if err := task.Add(DB); err != nil {
		respondErr(w, r, err, http.StatusBadRequest)
		return
	}

	respond(w, r, nil, http.StatusCreated)
}

func handleTaskGet(w http.ResponseWriter, r *http.Request) {
	params := parseParameters(r, "/task/int64:id")

	taskID, ok := params["id"]
	if !ok {
		respondErr(w, r, errors.New("incorrect task id"), http.StatusBadRequest)
		return
	}

	task, err := GetTask(taskID.(int64), DB)
	if err != nil {
		respondErr(w, r, err, http.StatusInternalServerError)
		return
	}
	respond(w, r, task, http.StatusOK)
}

func handleTaskDelete(w http.ResponseWriter, r *http.Request) {

}
