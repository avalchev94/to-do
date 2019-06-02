package web

import (
	"encoding/json"
	"net/http"

	"github.com/avalchev94/to-do-app/api"

	"github.com/gin-gonic/gin"
)

func dashboard(ctx *gin.Context) {
	templateData := gin.H{
		"Dashboard": true,
	}

	r, err := apiGET("/boards", ctx.Request.Cookies())
	if err != nil {
		templateData["Error"] = "Can't load dashboard... Try again?"
		ctx.Error(err)
		return
	}

	switch r.StatusCode {
	case http.StatusOK:
		var boards []gin.H
		// NOTE: not error handling
		json.NewDecoder(r.Body).Decode(&boards)
		templateData["Boards"] = boards
	case http.StatusUnauthorized:
		ctx.Redirect(http.StatusTemporaryRedirect, "/login")
	case http.StatusBadRequest:
		err := apiError(r)
		switch err {
		case api.ResourceNotFound:
			templateData["Error"] = "No dashboards found."
		default:
			templateData["Error"] = "Can't load dashboard... Try again?"
			ctx.Error(err.Error())
		}
	}

	ctx.HTML(http.StatusOK, "layout.html", templateData)
}
