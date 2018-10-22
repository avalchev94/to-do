package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func dashboard(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "layout.html", gin.H{
		"Dashboard": true,
	})
}
