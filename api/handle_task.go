package api

import (
	"encoding/json"
	"net/http"

	"github.com/avalchev94/to-do-app/database"
	"github.com/gin-gonic/gin"
)

func getTask(ctx *gin.Context) {
	task, err := database.GetTask(ctx.GetInt64("id"), db)
	handleErrorGet(task, err, ctx)
}

func postScheduledTask(ctx *gin.Context) {
	var t database.ScheduledTask
	if err := json.NewDecoder(ctx.Request.Body).Decode(&t); err != nil {
		ctx.JSON(http.StatusBadRequest, NewError(err))
		return
	}

	t.Task.UserID = ctx.GetInt64("id")
	handleErrorPost(t.Add(db), ctx)
}

func postRepetitiveTask(ctx *gin.Context) {
	var t database.RepetitiveTask
	if err := json.NewDecoder(ctx.Request.Body).Decode(&t); err != nil {
		ctx.JSON(http.StatusBadRequest, NewError(err))
		return
	}

	t.Task.UserID = ctx.GetInt64("id")
	handleErrorPost(t.Add(db), ctx)
}

func getScheduledTasks(ctx *gin.Context) {
	tasks, err := database.GetScheduledTasks(ctx.GetInt64("id"), db)
	handleErrorGet(tasks, err, ctx)
}

func getRepetitiveTasks(ctx *gin.Context) {
	// userID := user.(*database.User).ID
	// tasks, err := database.GetRepetitiveTasks(userID, db)
	// if !ok {
	// 	ctx.JSON(http.StatusInternalServerError, Error{err.Error()})
	// 	return
	// }
	// ctx.JSON(http.StatusOK, tasks)
}
