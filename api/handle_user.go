package main

import (
	"encoding/json"
	"errors"
	"net/http"
)

func handleUser(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		params := parseParameters(r, "/user/int64:id/string:relation")

		userID, ok := params["id"]
		if !ok {
			respondErr(w, r, errors.New("incorrect user id"), http.StatusBadRequest)
			return
		}

		relation, ok := params["relation"]
		if !ok {
			handleUserGet(w, r, userID.(int64))
			return
		}

		switch relation {
		case "labels":
			handleUserLabelsGet(w, r, userID.(int64))
		case "tasks":
			handleUserTasksGet(w, r, userID.(int64))
		}
	case "POST":
		handleUserRegister(w, r)
	case "DELETE":
		handleUserDelete(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func handleUserGet(w http.ResponseWriter, r *http.Request, userID int64) {
	user, err := GetUser(userID, DB)
	if err != nil {
		respondErr(w, r, err, http.StatusBadRequest)
		return
	}
	respond(w, r, user, http.StatusOK)
}

func handleUserLabelsGet(w http.ResponseWriter, r *http.Request, userID int64) {
	labels, err := GetLabels(userID, DB)
	if err != nil {
		respondErr(w, r, err, http.StatusBadRequest)
		return
	}
	respond(w, r, labels, http.StatusOK)
}

func handleUserTasksGet(w http.ResponseWriter, r *http.Request, userID int64) {
	tasks, err := GetTasks(userID, DB)
	if err != nil {
		respondErr(w, r, err, http.StatusBadRequest)
		return
	}
	respond(w, r, tasks, http.StatusOK)
}

func handleUserRegister(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		respondErr(w, r, err, http.StatusBadRequest)
		return
	}
	if err := user.Add(DB); err != nil {
		//TODO: Not always the error is InternalServerError
		respondErr(w, r, err, http.StatusInternalServerError)
		return
	}

	respond(w, r, nil, http.StatusCreated)
}

func handleUserDelete(w http.ResponseWriter, r *http.Request) {

}
