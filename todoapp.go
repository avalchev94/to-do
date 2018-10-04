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

func RunTaskRepeater(threads int) error {
	repeater, err := api.NewTaskRepeater()
	if err != nil {
		return err
	}

	log.Printf("Running TaskRepeater with %d threads\n", threads)
	go repeater.Run(threads)
	return nil
}
