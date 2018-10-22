package main

import (
	"flag"
	"log"
	"os"

	"github.com/go-redis/redis"

	"github.com/avalchev94/to-do-app"
)

func main() {
	flag.Parse()

	if flag.NArg() == 0 {
		log.Fatal("no arguments")
	}

	switch flag.Arg(0) {
	case "api":
		runAPI(flag.Args()[1:])
	case "web":
		runWeb(flag.Args()[1:])
	case "scheduler":
		runScheduler(flag.Args()[1:])
	}
}

func runAPI(args []string) {
	fs := flag.NewFlagSet("API", flag.ContinueOnError)
	httpAddr := fs.String("http", ":8080", "API http address")
	fs.Parse(args)

	err := todoapp.RunAPIServer(*httpAddr,
		todoapp.DBOptions{
			Host:     os.Getenv("PG_HOST"),
			Port:     os.Getenv("PG_PORT"),
			User:     os.Getenv("PG_USER"),
			Password: os.Getenv("PG_PASSWORD"),
			DBName:   "todo_app",
		},
		redis.Options{
			Addr:     os.Getenv("REDIS_ADDR"),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       0,
		})
	if err != nil {
		log.Fatal(err)
	}
}

func runWeb(args []string) {
	fs := flag.NewFlagSet("WEB", flag.ContinueOnError)
	apiAddr := fs.String("api", "http://localhost:8080", "API http address")
	httpAddr := fs.String("http", ":8081", "web http address")
	fs.Parse(args)

	if err := todoapp.RunWebClient(*httpAddr, *apiAddr); err != nil {
		log.Fatal(err)
	}
}

func runScheduler(args []string) {
	fs := flag.NewFlagSet("Scheduler", flag.ContinueOnError)
	threads := fs.Int("threads", 10, "Threads used for repetitive task creation")
	fs.Parse(args)

	err := todoapp.RunTaskScheduler(*threads,
		todoapp.DBOptions{
			Host:     os.Getenv("PG_HOST"),
			Port:     os.Getenv("PG_PORT"),
			User:     os.Getenv("PG_USER"),
			Password: os.Getenv("PG_PASSWORD"),
			DBName:   "todo_app",
		},
		redis.Options{
			Addr:     os.Getenv("REDIS_ADDR"),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       0,
		})
	if err != nil {
		log.Fatalln(err)
	}
}
