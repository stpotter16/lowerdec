package main

import (
    "embed"
    "fmt"
    "html/template"
    "io/fs"
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
        fmt.Fprintf(w, "Welcome!")
    })
    http.HandleFunc("/static/style.css", func (w http.ResponseWriter, r *http.Request) {
        http.FileServer(http.FS(css))
    })

    http.ListenAndServe(":8080", nil)
}

