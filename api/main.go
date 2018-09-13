package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/go-redis/redis"
	// usign postgres driver
	_ "github.com/lib/pq"
)

// DB is connection to the database.
var DB *sql.DB

// RedisConn is connection to the redis service.
var RedisConn *redis.Client

func main() {
	var addr = flag.String("address", ":8080", "the address of the api")
	flag.Parse()

	// Establish PostgresSQL connection
	var err error
	DB, err = ConnectDB(ConnectionInfo{
		Host:     os.Getenv("PG_HOST"),
		Port:     os.Getenv("PG_PORT"),
		User:     os.Getenv("PG_USER"),
		Password: os.Getenv("PG_PASSWORD"),
		DBName:   "todo_app",
	})
	if err != nil {
		log.Fatalln(err)
	}
	defer DB.Close()

	// Establish Redis connection
	RedisConn = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
	if RedisConn == nil {
		log.Fatalln("redis connection failed")
	}
	defer RedisConn.Close()

	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/user/", handleUser)
	http.HandleFunc("/task/", handleTask)
	http.HandleFunc("/label/", handleLabel)

	if err = http.ListenAndServe(*addr, nil); err != nil {
		log.Fatalln(err)
	}
}
