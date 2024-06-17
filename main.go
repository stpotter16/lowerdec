package main

import (
	"context"
	_ "database/sql"
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/mattn/go-sqlite3"
)

var (
	//go:embed  static/css/output.css
	css embed.FS

	//go:embed all:templates/*
	templatesFS embed.FS

	//go:embed static/assets/favicon.ico
	favicon embed.FS
	//go:embed static/assets/logo.png
	navLogo embed.FS
	//go:embed static/assets/logo-light.svg
	svgLog embed.FS
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

	_, err := initDB()
	if err != nil {
		return err
	}

	html, err := template.New("").ParseFS(templatesFS, "templates/*.html")
	if err != nil {
		return err
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Host: %s", r.Host)
		log.Printf("%s: %s", r.Method, r.URL.Path)
		html.ExecuteTemplate(w, "landing.html", nil)
	})
	http.HandleFunc("/start", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s: %s", r.Method, r.URL.Path)
		if r.Method != http.MethodPost {
			html.ExecuteTemplate(w, "start.html", nil)
			return
		}
		log.Printf("%s: %s", "dec", r.FormValue("dec"))
		_, handler, _ := r.FormFile("dec")
		log.Printf("File size: %d", handler.Size)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		html.ExecuteTemplate(w, "start.html", struct{ Success bool }{true})
	})

	http.HandleFunc("agent.lowerdec.localhost/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Host: %s", r.Host)
		log.Printf("%s: %s", r.Method, r.URL.Path)
		fmt.Fprintf(w, "Hello from agents")
	})

	http.HandleFunc("/static/css/output.css", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s: %s", r.Method, r.URL.Path)
		handler := http.FileServer(http.FS(css))
		handler.ServeHTTP(w, r)
	})
	http.HandleFunc("/static/assets/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s: %s", r.Method, r.URL.Path)
		handler := http.FileServer(http.FS(favicon))
		handler.ServeHTTP(w, r)
	})
	http.HandleFunc("/static/assets/logo.png", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s: %s", r.Method, r.URL.Path)
		handler := http.FileServer(http.FS(navLogo))
		handler.ServeHTTP(w, r)
	})
	http.HandleFunc("/static/assets/logo-light.svg", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s: %s", r.Method, r.URL.Path)
		handler := http.FileServer(http.FS(svgLog))
		handler.ServeHTTP(w, r)
	})

	http.ListenAndServe(":8080", nil)

	<-ctx.Done()
	log.Print("Received termination signal. Shutting down")

	return nil
}
