package main

import (
	"flag"
	"log"

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
	}
}

func runAPI(args []string) {
	fs := flag.NewFlagSet("API", flag.ContinueOnError)
	httpAddr := fs.String("http", ":8080", "API http address")
	threads := fs.Int("threads", 10, "Threads used for repetitive task creation")
	fs.Parse(args)

	if err := todoapp.RunTaskRepeater(*threads); err != nil {
		log.Fatal(err)
	}

	if err := todoapp.RunAPI(*httpAddr); err != nil {
		log.Fatal(err)
	}
}
