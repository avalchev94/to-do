package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/avalchev94/to-do-app/database"
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
	var task database.ScheduledTask
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		respondErr(w, r, err, http.StatusBadRequest)
		return
	}

	var err error
	task.Task.UserID, err = getCurrentUserID(r)
	if err != nil {
		respondErr(w, r, errors.New("no logged user"), http.StatusBadRequest)
		return
	}

	if err := task.Add(db); err != nil {
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

	task, err := database.GetScheduledTask(taskID.(int64), db)
	if err != nil {
		respondErr(w, r, err, http.StatusBadRequest)
		return
	}

	respond(w, r, task, http.StatusOK)
}

func handleTaskDelete(w http.ResponseWriter, r *http.Request) {
}
