package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/avalchev94/to-do-app/api"
)

var (
	apiAddress string
)

func apiGET(url string, cookies []*http.Cookie) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, apiAddress+url, nil)
	if err != nil {
		return nil, err
	}

	if cookies != nil {
		for _, c := range cookies {
			req.AddCookie(c)
		}
	}
	return http.DefaultClient.Do(req)
}

func apiPOST(url string, data interface{}, cookies []*http.Cookie) (*http.Response, error) {
	b, err := json.Marshal(data)
	fmt.Println(string(b), "here")
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, apiAddress+url, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}

	if cookies != nil {
		for _, c := range cookies {
			req.AddCookie(c)
		}
	}
	return http.DefaultClient.Do(req)
}

func apiError(r *http.Response) api.Error {
	var apiError api.Error
	if err := json.NewDecoder(r.Body).Decode(&apiError); err != nil {
		return api.NewError(err)
	}
	return apiError
}
