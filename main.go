package main

import (
    "context"
    "database/sql"
    "embed"
    "fmt"
    "html/template"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"

    _ "github.com/mattn/go-sqlite3"
    "github.com/google/uuid"
)

var (
    //go:embed  static/style.css
    css embed.FS

    //go:embed all:templates/*
    templatesFS embed.FS

    //go:embed static/assets/favicon.ico
    favicon embed.FS
    //go:embed static/assets/logo.png
    navLogo embed.FS
)

const DB_PATH = "DB_PATH"

func run() error {
    ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM)
    defer stop()

    dbPath, present := os.LookupEnv(DB_PATH)

    if !present {
        return fmt.Errorf("%s not set", DB_PATH)
    }

    db, err := sql.Open("sqlite3", dbPath)
    if err != nil {
        return err
    }
    defer db.Close()

    if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS signups (userId STRING PRIMARY KEY, name STRING, email STRING);`); err != nil {
        return fmt.Errorf("Cannot create table: %w", err)
    }

    html, err := template.New("").ParseFS(templatesFS, "templates/*.html")
    if err != nil {
        return err
    }

    http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
        log.Printf("%s: %s", r.Method, r.URL.Path)
        html.ExecuteTemplate(w, "landing.html", nil)
    })
    http.HandleFunc("/start", func (w http.ResponseWriter, r *http.Request) {
        log.Printf("%s: %s", r.Method, r.URL.Path)
        if r.Method != http.MethodPost {
            html.ExecuteTemplate(w, "start.html", nil)
            return
        }
        log.Printf("%s: %s", "full name", r.FormValue("name"))
        log.Printf("%s: %s", "email", r.FormValue("email"))
        userId := uuid.New()
        if _, err := db.Exec(`INSERT INTO signups (userId, name, email) VALUES (?, ?, ?)`, userId.String(), r.FormValue("name"), r.FormValue("email")); err != nil {
            // TODO - should return an error
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
        html.ExecuteTemplate(w, "start.html", struct{ Success bool }{true})
    })
    http.HandleFunc("/static/style.css", func (w http.ResponseWriter, r *http.Request) {
        log.Printf("%s: %s", r.Method, r.URL.Path)
        handler := http.FileServer(http.FS(css))
        handler.ServeHTTP(w, r)
    })
    http.HandleFunc("/static/assets/favicon.ico", func (w http.ResponseWriter, r *http.Request) {
        log.Printf("%s: %s", r.Method, r.URL.Path)
        handler := http.FileServer(http.FS(favicon))
        handler.ServeHTTP(w, r)
    })
    http.HandleFunc("/static/assets/logo.png", func (w http.ResponseWriter, r *http.Request) {
        log.Printf("%s: %s", r.Method, r.URL.Path)
        handler := http.FileServer(http.FS(navLogo))
        handler.ServeHTTP(w, r)
    })

    http.ListenAndServe(":8080", nil)

    <-ctx.Done()
    log.Print("Received termination signal. Shutting down")

    return nil
}

func main() {
    if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

