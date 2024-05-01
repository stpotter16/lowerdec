package main

import (
    "embed"
    "html/template"
    "io/fs"
    "log"
    "net/http"
)

var (
    //go:embed  static/style.css
    css embed.FS

    //go:embed all:templates/*
    templatesFS embed.FS

    //go:embed static/assets/favicon.ico
    favicon embed.FS

    html *template.Template
)

func parseTemplates(templates fs.FS) *template.Template {
    parsed := template.Must(template.New("").ParseFS(templates, "templates/*.html"))
    return parsed
}

func main() {
    html = parseTemplates(templatesFS)

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
        html.ExecuteTemplate(w, "start.html", struct{ Success bool }{true})
    })
    http.HandleFunc("/static/style.css", func (w http.ResponseWriter, r *http.Request) {
        log.Printf("%s: %s", r.Method, r.URL.Path)
        handler := http.FileServer(http.FS(css))
        handler.ServeHTTP(w, r)
    })
    http.HandleFunc("/static/assets/favicon.ico", func (w http.ResponseWriter, r *http.Request) {
        log.Printf("In the favicon handler")
        log.Printf("%s: %s", r.Method, r.URL.Path)
        handler := http.FileServer(http.FS(favicon))
        handler.ServeHTTP(w, r)
    })

    log.Fatal(http.ListenAndServe(":8080", nil))
}

