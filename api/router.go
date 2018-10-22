package api

import (
	"database/sql"

	authenticator "github.com/avalchev94/to-do-app/api/auth"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

var (
	router *gin.Engine
	auth   *authenticator.Auth
	db     *sql.DB
)

// Router returns gin.Engine with all the routes API need.
// All you need to do is call: http.ListenAndServer(":PORT", Router(db, redisClient)).
// Where:
// - db is sql.DB with Postgres driver
// - redisClient is self-descriptive
func Router(sqlDB *sql.DB, redisClient *redis.Client) *gin.Engine {
	if router != nil {
		return router
	}

	db = sqlDB
	auth = authenticator.New(redisClient, "AuthCookie")
	router = gin.Default()
	router.GET("/login", getLogged, getLogin)
	router.POST("/login", postLogin)
	router.GET("/logout", logout)

	router.GET("/user/:id", getID, getUser)
	router.GET("/user/:id/labels", getID, getUserLabels)
	router.GET("/user/:id/scheduled_tasks", getID, getScheduledTasks)
	router.GET("/user/:id/repetitive_tasks", getID, getRepetitiveTasks)
	router.POST("/user", postUser)

	router.GET("/label/:id", getID, getLabel)
	router.GET("/labels", getLogged, getUserLabels)
	router.POST("/label", getLogged, postLabel)

	router.GET("/task/:id", getID, getTask)
	router.GET("/scheduled_tasks", getLogged, getScheduledTasks)
	router.GET("/repetitive_tasks", getLogged, getRepetitiveTasks)
	router.POST("/scheduled_task", getLogged, postScheduledTask)
	router.POST("/repetitive_task", getLogged, postRepetitiveTask)

	return router
}
