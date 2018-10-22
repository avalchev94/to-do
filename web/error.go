package web

import (
	"github.com/gin-gonic/gin"
)

func newError(title, description string) gin.H {
	return gin.H{
		"Title":       title,
		"Description": description,
	}
}
