package api

import (
	"encoding/json"
	"net/http"

	"github.com/avalchev94/to-do-app/database"
	"github.com/gin-gonic/gin"
)

func getLogin(ctx *gin.Context) {
	user, ok := ctx.Get("user")
	if !ok {
		ctx.Status(http.StatusInternalServerError)
		return
	}
	ctx.JSON(http.StatusOK, user)
}

func postLogin(ctx *gin.Context) {
	if _, err := auth.Logged(ctx); err == nil {
		ctx.JSON(http.StatusBadRequest, AlreardyAuth)
		return
	}

	loginData := struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}{}
	if err := json.NewDecoder(ctx.Request.Body).Decode(&loginData); err != nil {
		ctx.JSON(http.StatusBadRequest, NewError(err))
		return
	}

	user, err := database.VerifyLogin(loginData.Name, loginData.Password, db)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, NewError(err))
		return
	}

	if err := auth.Login(user.ID, ctx); err != nil {
		ctx.Status(http.StatusInternalServerError)
		ctx.Error(err)
		return
	}

	ctx.Status(http.StatusOK)
}

func logout(ctx *gin.Context) {
	if err := auth.Logout(ctx); err != nil {
		ctx.JSON(http.StatusBadRequest, NewError(err))
		return
	}
	ctx.Status(http.StatusOK)
}
