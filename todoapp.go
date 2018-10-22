package todoapp

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/avalchev94/to-do-app/api"
	"github.com/avalchev94/to-do-app/api/scheduler"
	"github.com/avalchev94/to-do-app/web"
	"github.com/go-redis/redis"

	// blank importing postgres driver
	_ "github.com/lib/pq"
)

// DBOptions wraps all the needed information for establishing connection with PostgresSQL.
type DBOptions struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func (db *DBOptions) connectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		db.Host, db.Port, db.User, db.Password, db.DBName)
}

func RunAPIServer(httpAddr string, dbOpt DBOptions, redisOpt redis.Options) error {
	db, err := sql.Open("postgres", dbOpt.connectionString())
	if err != nil {
		return err
	}

	redisClient := redis.NewClient(&redisOpt)
	if err := redisClient.Ping().Err(); err != nil {
		return err
	}

	return http.ListenAndServe(httpAddr, api.Router(db, redisClient))
}

func RunWebClient(httpAddr string, apiAddr string) error {
	return http.ListenAndServe(httpAddr, web.Router(apiAddr))
}

func RunTaskScheduler(threads int, dbOpt DBOptions, redisOpt redis.Options) error {
	db, err := sql.Open("postgres", dbOpt.connectionString())
	if err != nil {
		return err
	}

	redisClient := redis.NewClient(&redisOpt)
	if err := redisClient.Ping().Err(); err != nil {
		return err
	}

	scheduler, err := scheduler.New(redisClient, db)
	if err != nil {
		return err
	}

	return scheduler.Run(threads)
}
