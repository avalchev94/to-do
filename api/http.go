package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type responseData struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func respond(w http.ResponseWriter, r *http.Request, v interface{}, code int) {
	response := responseData{
		Success: true,
		Data:    v,
	}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(response)
	if err != nil {
		respondErr(w, r, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_, err = buf.WriteTo(w)
	if err != nil {
		log.Printf("respond: %s", err)
	}
}

func respondErr(w http.ResponseWriter, r *http.Request, err error, code int) {
	response := responseData{
		Success: false,
		Error:   err.Error(),
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("respondErr: %s", err)
	}
}

func parseParameters(r *http.Request, pattern string) map[string]interface{} {
	params := make(map[string]interface{})
	pathSegs := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	for i, seg := range strings.Split(strings.Trim(pattern, "/"), "/") {
		if i > len(pathSegs)-1 {
			return params
		}

		pair := strings.Split(seg, ":")
		if len(pair) < 2 {
			pair = append([]string{"string"}, pair...)
		}

		switch pair[0] {
		case "int64":
			int, err := strconv.ParseInt(pathSegs[i], 10, 64)
			if err != nil {
				continue
			}
			params[pair[1]] = int
		default:
			params[pair[1]] = pathSegs[i]
		}
	}
	return params
}

func getCurrentUserID(r *http.Request) (int64, error) {
	cookie, err := r.Cookie(AuthCookie)
	if err != nil {
		return 0, err
	}
	return GetSessionUser(cookie.Value)
}
