package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func root(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "layout.html", nil)
}
