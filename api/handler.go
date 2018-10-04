package api

import (
	"net/http"
)

// Handler returns http.ServeMux that handles every pattern of the API.
func Handler() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/login", handleLogin)
	mux.HandleFunc("/user/", handleUser)
	mux.HandleFunc("/task/", handleTask)
	mux.HandleFunc("/label/", handleLabel)
	mux.HandleFunc("/repetitive_task/", handleRepetitiveTask)

	return mux
}
