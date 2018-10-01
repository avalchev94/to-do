package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/avalchev94/to-do-app/database"
)

func handleRepetitiveTask(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		handleRepetitiveTaskGet(w, r)
	case "POST":
		handleRepetitiveTaskPost(w, r)
	case "DELETE":
		handleRepetitiveTaskDelete(w, r)
	}
}

type RepetitiveTask struct {
	Task       *database.Task           `json:"task"`
	Repeatance *database.TaskRepeatance `json:"repeatance"`
}

func handleRepetitiveTaskGet(w http.ResponseWriter, r *http.Request) {
	params := parseParameters(r, "/repetitive_task/int64:id")

	taskID, ok := params["id"]
	if !ok {
		respondErr(w, r, errors.New("incorrect task id"), http.StatusBadRequest)
		return
	}

	task, repeatance, err := database.GetTaskAndRepeatance(taskID.(int64), db)
	if err != nil {
		respondErr(w, r, err, http.StatusBadRequest)
		return
	}
	respond(w, r, RepetitiveTask{task, repeatance}, http.StatusOK)
}

func handleRepetitiveTaskPost(w http.ResponseWriter, r *http.Request) {
	var t RepetitiveTask
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		respondErr(w, r, err, http.StatusBadRequest)
		return
	}

	var err error
	t.Task.UserID, err = getCurrentUserID(r)
	if err != nil {
		respondErr(w, r, errors.New("no logged user"), http.StatusBadRequest)
		return
	}

	if err := database.AddRepetitiveTask(t.Task, t.Repeatance, db); err != nil {
		respondErr(w, r, err, http.StatusBadRequest)
		return
	}
	respond(w, r, nil, http.StatusCreated)
}

func handleRepetitiveTaskDelete(w http.ResponseWriter, r *http.Request) {
}
