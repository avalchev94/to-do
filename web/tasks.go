package web

import (
	"encoding/json"
	"net/http"

	"github.com/avalchev94/to-do-app/api"

	"github.com/gin-gonic/gin"
)

func handleTasks(ctx *gin.Context) {
	r, err := apiGET("/scheduled_tasks", ctx.Request.Cookies())
	if err != nil {
		ctx.HTML(http.StatusOK, "task.html", gin.H{
			"Error": "Can't load tasks...",
		})
		ctx.Error(err)
		return
	}

	switch r.StatusCode {
	case http.StatusOK:
		var tasks []map[string]interface{}
		// NOTE: not error handling
		json.NewDecoder(r.Body).Decode(&tasks)
		ctx.HTML(http.StatusOK, "tasks.html", gin.H{
			"Tasks": tasks,
		})
	case http.StatusUnauthorized:
		ctx.Redirect(http.StatusTemporaryRedirect, "/login")
	case http.StatusBadRequest:
		err := apiError(r)
		if err == api.ResourceNotFound {
			ctx.HTML(http.StatusOK, "tasks.html", gin.H{
				"Error": "No tasks added...",
			})
		}
	}

	if !ctx.Writer.Written() {
		ctx.HTML(http.StatusOK, "tasks.html", gin.H{
			"Error": "Can't load tasks...",
		})
	}
}

func postTask(ctx *gin.Context) {
	task := map[string]interface{}{}
	json.NewDecoder(ctx.Request.Body).Decode(&task)

	var url string
	if _, ok := task["repeatance"]; ok {
		url = "/repetitive_task"
	} else {
		url = "/scheduled_task"
	}

	r, err := apiPOST(url, task, ctx.Request.Cookies())
	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	switch r.StatusCode {
	case http.StatusCreated:
		ctx.Status(http.StatusCreated)
	case http.StatusBadRequest:
		ctx.JSON(http.StatusBadRequest, apiError(r))
	default:
		ctx.Status(http.StatusInternalServerError)
	}
}

func getLabels(ctx *gin.Context) {
	r, err := apiGET("/labels", ctx.Request.Cookies())
	if err != nil {
		ctx.Error(err)
		ctx.Status(http.StatusInternalServerError)
	}

	switch r.StatusCode {
	case http.StatusOK, http.StatusBadRequest:
		var data interface{}
		json.NewDecoder(r.Body).Decode(&data)
		ctx.JSON(r.StatusCode, data)
	default:
		ctx.Status(http.StatusInternalServerError)
	}
}
