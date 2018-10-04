package api

import (
	"database/sql"
	"log"
	"os"

	"github.com/avalchev94/to-do-app/database"
	"github.com/go-redis/redis"
)

var (
	db        *sql.DB
	redisConn *redis.Client
)

func init() {
	var err error
	db, err = database.Connect(database.ConnectionInfo{
		Host:     os.Getenv("PG_HOST"),
		Port:     os.Getenv("PG_PORT"),
		User:     os.Getenv("PG_USER"),
		Password: os.Getenv("PG_PASSWORD"),
		DBName:   "todo_app",
	})
	if err != nil {
		log.Println("API package: ", err)
	}
	redisConn = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
	if redisConn == nil {
		log.Println("API package: redis connection failed")
	}
}
