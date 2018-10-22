package api

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/avalchev94/to-do-app/database"
)

func postLabel(ctx *gin.Context) {
	var label database.Label
	if err := json.NewDecoder(ctx.Request.Body).Decode(&label); err != nil {
		ctx.JSON(http.StatusBadRequest, NewError(err))
		return
	}

	label.ID = ctx.GetInt64("id")
	handleErrorPost(label.Add(db), ctx)
}

func getLabel(ctx *gin.Context) {
	label, err := database.GetLabel(ctx.GetInt64("id"), db)
	handleErrorGet(label, err, ctx)
}
