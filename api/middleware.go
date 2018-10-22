package api

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"

	"github.com/avalchev94/to-do-app/database"
	"github.com/gin-gonic/gin"
)

func getID(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, IncorrectParameter)
		return
	}
	ctx.Set("id", id)
	ctx.Next()
}

func getLogged(ctx *gin.Context) {
	userID, err := auth.Logged(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Unauthorized)
		return
	}

	user, err := database.GetUser(userID, db)
	if err != nil {
		switch auth.Logout(ctx) {
		case nil:
			ctx.JSON(http.StatusUnauthorized, Unauthorized)
		default:
			ctx.JSON(http.StatusInternalServerError, nil)
		}
		return
	}

	ctx.Set("user", user)
	ctx.Set("id", user.ID)
	ctx.Next()
}

func handleErrorGet(data interface{}, err error, ctx *gin.Context) {
	switch err {
	case nil:
		ctx.JSON(http.StatusOK, data)
	case sql.ErrNoRows:
		ctx.JSON(http.StatusBadRequest, ResourceNotFound)
	default:
		ctx.Status(http.StatusInternalServerError)
		ctx.Error(err)
	}
}

func handleErrorPost(err error, ctx *gin.Context) {
	switch {
	case err == nil:
		ctx.Status(http.StatusCreated)
	case strings.HasPrefix(err.Error(), "db:"):
		ctx.JSON(http.StatusBadRequest, NewError(err))
	default:
		ctx.Status(http.StatusInternalServerError)
		ctx.Error(err)
	}
}
