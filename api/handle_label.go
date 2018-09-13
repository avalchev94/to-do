package main

import (
	"encoding/json"
	"errors"
	"net/http"
)

func handleLabel(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		handlePostLabel(w, r)
	case "GET":
		handleGetLabel(w, r)
	case "DELETE":
		handleDeleteLabel(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func handlePostLabel(w http.ResponseWriter, r *http.Request) {
	var label Label
	if err := json.NewDecoder(r.Body).Decode(&label); err != nil {
		respondErr(w, r, err, http.StatusBadRequest)
		return
	}

	// Get the logged user
	userID, err := getCurrentUserID(r)
	if err != nil {
		respondErr(w, r, errors.New("no logged user"), http.StatusBadRequest)
		return
	}
	label.UserID = userID

	if err := label.Add(DB); err != nil {
		//TODO: Not always the error is InternalServerError
		respondErr(w, r, err, http.StatusInternalServerError)
		return
	}

	respond(w, r, nil, http.StatusCreated)
}

func handleGetLabel(w http.ResponseWriter, r *http.Request) {
	params := parseParameters(r, "/label/int64:id")

	userID, ok := params["id"]
	if !ok {
		respondErr(w, r, errors.New("incorrect label id"), http.StatusBadRequest)
		return
	}

	label, err := GetLabel(userID.(int64), DB)
	if err != nil {
		respondErr(w, r, err, http.StatusBadRequest)
		return
	}
	respond(w, r, &label, http.StatusOK)
}

func handleDeleteLabel(w http.ResponseWriter, r *http.Request) {

}
