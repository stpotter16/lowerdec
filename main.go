package main

import (
    "context"
    "database/sql"
    "embed"
    "fmt"
    "html/template"
    "log"
    "net/http"
    "net/url"
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

func check_db_settings(db *sql.DB) error {
    busy_timeout_row := db.QueryRow("PRAGMA busy_timeout")
    if busy_timeout_row == nil {
        return fmt.Errorf("PRAMA busy_timeout not found")
    }
    var busy_timeout int
    if err := busy_timeout_row.Scan(&busy_timeout); err != nil {
        return err
    }
    log.Printf("Busy timeout set to %d", busy_timeout)

    sync_mode_row := db.QueryRow("PRAGMA synchronous")
    if sync_mode_row == nil {
        return fmt.Errorf("PRAGMA synchronous not found")
    }
    var sync_mode int
    if err := sync_mode_row.Scan(&sync_mode); err != nil {
        return err
    }
    log.Printf("Synchronous mode set to %d", sync_mode)

    journal_mode_row := db.QueryRow("PRAGMA journal_mode")
    if journal_mode_row == nil {
        return fmt.Errorf("PRAMA journal_mode not found")
    }
    var journal_mode string
    if err := journal_mode_row.Scan(&journal_mode); err != nil {
        return err
    }
    log.Printf("Journal mode set to %s", journal_mode)

    return nil
}

func run() error {
    ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM)
    defer stop()

    dbPath, present := os.LookupEnv(DB_PATH)

    if !present {
        return fmt.Errorf("%s not set", DB_PATH)
    }

    // TODO - This should really be configured via environment variables
    options := url.QueryEscape("_timeout=5000&_sync=1")

    dsn := "file:" + dbPath + "?" + options
    log.Printf("%s", dsn)

    db, err := sql.Open("sqlite3", dsn)
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

    if err := check_db_settings(db); err != nil {
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
        log.Printf("%s: %s", "dec", r.FormValue("dec"))
        _ , handler, _ := r.FormFile("dec")
        log.Printf("File size: %d", handler.Size)
        if err != nil {
            fmt.Fprintln(os.Stderr, err)
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
        html.ExecuteTemplate(w, "start.html", struct{ Success bool }{true})
    })
    http.HandleFunc("/static/css/output.css", func (w http.ResponseWriter, r *http.Request) {
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
    http.HandleFunc("/static/assets/logo-light.svg", func (w http.ResponseWriter, r *http.Request) {
        log.Printf("%s: %s", r.Method, r.URL.Path)
        handler := http.FileServer(http.FS(svgLog))
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

