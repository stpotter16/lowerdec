package main

import (
	"context"
	_ "database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/mattn/go-sqlite3"
)

const DB_PATH = "DB_PATH"

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM)
	defer stop()

	err := initDB()
	if err != nil {
		return err
	}

	err = parseHtml()
	if err != nil {
		return err
	}

    mux := http.NewServeMux()

	// Users
	mux.HandleFunc("GET /", LowerDecLanding)
	mux.HandleFunc("GET /start", GetPolicyStart)
	mux.HandleFunc("POST /start", PostPolicyStart)

	// Agents
	mux.HandleFunc("GET agent.lowerdec.localhost/", AgengLowerDecLanding)

	// Common
	mux.Handle("GET /static/", http.FileServer(http.FS(staticFS)))

    handler := Logging(mux)

	http.ListenAndServe(":8080", handler)

	<-ctx.Done()
	log.Print("Received termination signal. Shutting down")

	return nil
}
