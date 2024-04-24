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
    http.HandleFunc("/about", func (w http.ResponseWriter, r *http.Request) {
        log.Printf("%s: %s", r.Method, r.URL.Path)
        html.ExecuteTemplate(w, "about.html", nil)
    })
    http.HandleFunc("/start", func (w http.ResponseWriter, r *http.Request) {
        log.Printf("%s: %s", r.Method, r.URL.Path)
        if r.Method != http.MethodPost {
            html.ExecuteTemplate(w, "start.html", nil)
            return
        }
        log.Printf("%s: %s", "full name:", r.FormValue("name"))
        log.Printf("%s: %s", "email", r.FormValue("email"))
    })
    http.HandleFunc("/static/style.css", func (w http.ResponseWriter, r *http.Request) {
        log.Printf("%s: %s", r.Method, r.URL.Path)
        handler := http.FileServer(http.FS(css))
        handler.ServeHTTP(w, r)
    })

    log.Fatal(http.ListenAndServe(":8080", nil))
}

