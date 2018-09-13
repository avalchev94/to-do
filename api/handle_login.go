package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/avalchev94/to-do-app/database"
)

const (
	authCookie = "AuthCookie"
)

func handleLogin(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		handleLoginGet(w, r)
	case "POST":
		handleLoginPost(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func handleLoginGet(w http.ResponseWriter, r *http.Request) {
	userID, err := getCurrentUserID(r)
	if err != nil {
		respondErr(w, r, errors.New("no logged user"), http.StatusBadRequest)
		return
	}

	user, err := database.GetUser(userID, db)
	if err != nil {
		respondErr(w, r, err, http.StatusBadRequest)
		return
	}
	respond(w, r, user, http.StatusOK)
}

func handleLoginPost(w http.ResponseWriter, r *http.Request) {
	if _, err := getCurrentUserID(r); err == nil {
		respondErr(w, r, errors.New("already logged"), http.StatusBadRequest)
		return
	}

	loginData := struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&loginData); err != nil {
		respondErr(w, r, err, http.StatusBadRequest)
		return
	}
	user, err := database.VerifyLogin(loginData.Name, loginData.Password, db)
	if err != nil {
		respondErr(w, r, err, http.StatusBadRequest)
		return
	}

	uuid, err := newLoginSession(user.ID)
	if err != nil {
		respondErr(w, r, err, http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:  authCookie,
		Value: uuid,
	})
	respond(w, r, nil, http.StatusOK)
}
