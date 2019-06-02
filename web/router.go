package web

import (
	"path/filepath"
	"runtime"

	"github.com/gin-gonic/gin"
)

var (
	router *gin.Engine
)

// Router returns an gin.Engine pointer with all the routes for the web server.
// All you need to do is call: http.ListenAndServer(":PORT", Router(apiAddr))
func Router(apiAddr string) *gin.Engine {
	if router != nil {
		return router
	}

	apiAddress = apiAddr

	router = gin.Default()
	router.LoadHTMLGlob(templatesPattern())
	router.Static("/static", staticFolder())

	router.GET("/dash", dashboard)
	router.GET("/login", getLogin)
	router.POST("/login", postLogin)
	router.GET("/logout", logout)
	router.POST("/task", postTask)
	router.GET("/labels", getLabels)

	return router
}

func templatesPattern() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return ""
	}
	return filepath.Dir(file) + "/templates/*"
}

func staticFolder() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return ""
	}
	return filepath.Dir(file) + "/static"
}
