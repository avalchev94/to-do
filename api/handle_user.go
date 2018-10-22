package api

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/avalchev94/to-do-app/database"
)

func getUser(ctx *gin.Context) {
	user, err := database.GetUser(ctx.GetInt64("id"), db)
	switch err {
	case sql.ErrNoRows:
		ctx.JSON(http.StatusBadRequest, ResourceNotFound)
	case nil:
		ctx.JSON(http.StatusOK, user)
	default:
		ctx.Status(http.StatusInternalServerError)
		ctx.Error(err)
	}
}

func getUserLabels(ctx *gin.Context) {
	labels, err := database.GetLabels(ctx.GetInt64("id"), db)
	handleErrorGet(labels, err, ctx)
}

func postUser(ctx *gin.Context) {
	if _, err := auth.Logged(ctx); err == nil {
		ctx.JSON(http.StatusBadRequest, AlreardyAuth)
		return
	}

	var user database.User
	if err := json.NewDecoder(ctx.Request.Body).Decode(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, NewError(err))
		return
	}

	handleErrorPost(user.Add(db), ctx)
}
