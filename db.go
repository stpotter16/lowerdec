package main

import (
    "database/sql"
    "fmt"
    "log"
    "net/url"
    "os"

    _ "github.com/mattn/go-sqlite3"
)

func initDB() (*sql.DB, error) {
    // Collapse all this DB stuff
    dbPath, present := os.LookupEnv(DB_PATH)

    if !present {
        return nil, fmt.Errorf("%s not set", DB_PATH)
    }

    // TODO - This should really be configured via environment variables
    options := url.QueryEscape("_timeout=5000&_sync=1")

    dsn := "file:" + dbPath + "?" + options
    log.Printf("%s", dsn)

    db, err := sql.Open("sqlite3", dsn)
    if err != nil {
        return nil, err
    }
    defer db.Close()

    if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS signups (userId STRING PRIMARY KEY, name STRING, email STRING);`); err != nil {
        return nil, fmt.Errorf("Cannot create table: %w", err)
    }

    if err := check_db_settings(db); err != nil {
        return nil, err
    }

    return db, nil
}

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

