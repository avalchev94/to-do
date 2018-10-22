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
