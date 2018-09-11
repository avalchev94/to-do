package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

const (
	//AuthCookie is the name of the Authentication Cookie
	AuthCookie = "AuthCookie"
)

func handleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		respondErr(w, r, err, http.StatusInternalServerError)
		return
	}

	if err := user.Add(DB); err != nil {
		respondErr(w, r, err, http.StatusInternalServerError)
		return
	}

	respond(w, r, nil, http.StatusCreated)
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if _, err := r.Cookie(AuthCookie); err == nil {
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
	user, err := VerifyLoginData(loginData.Name, loginData.Password, DB)
	if err != nil {
		respondErr(w, r, err, http.StatusBadRequest)
		return
	}

	uuid, err := NewLoginSession(user.ID)
	if err != nil {
		respondErr(w, r, err, http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:  AuthCookie,
		Value: uuid,
	})
	respond(w, r, nil, http.StatusOK)
}

func handleUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	params := pathParams(r, "/user/:id")
	paramID, ok := params[":id"]

	var userID int64
	if ok {
		var err error
		userID, err = strconv.ParseInt(paramID, 10, 64)
		if err != nil {
			respondErr(w, r, fmt.Errorf("%s is not an id", paramID), http.StatusBadRequest)
			return
		}
	} else {
		cookie, err := r.Cookie(AuthCookie)
		if err != nil {
			respondErr(w, r, err, http.StatusBadRequest)
			return
		}
		if userID, err = GetSessionUser(cookie.Value); err != nil {
			respondErr(w, r, err, http.StatusBadRequest)
			return
		}
	}

	user, err := GetUser(userID, DB)
	if err != nil {
		respondErr(w, r, err, http.StatusBadRequest)
		return
	}

	respond(w, r, user, http.StatusOK)
}
