package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/avalchev94/to-do-app/database"
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
	var label database.Label
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

	if err := label.Add(db); err != nil {
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

	label, err := database.GetLabel(userID.(int64), db)
	if err != nil {
		respondErr(w, r, err, http.StatusBadRequest)
		return
	}
	respond(w, r, &label, http.StatusOK)
}

func handleDeleteLabel(w http.ResponseWriter, r *http.Request) {

}
