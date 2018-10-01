package todoapp

import (
	"log"
	"net/http"

	"github.com/avalchev94/to-do-app/api"
)

func RunAPI(httpAddr string) error {
	log.Println("Running API on", httpAddr)
	return http.ListenAndServe(httpAddr, api.Handler())
}
