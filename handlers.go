package main

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
)

var (
	//go:embed all:templates/*
	templatesFS embed.FS
	//go:embed static
	staticFS embed.FS
)

var html *template.Template

func parseHtml() error {
	var err error
	html, err = template.New("").ParseFS(templatesFS, "templates/*.html")
	if err != nil {
		return err
	}
	return nil
}

// Users
func LowerDecLanding(w http.ResponseWriter, r *http.Request) {
	html.ExecuteTemplate(w, "landing.html", nil)
}

func GetPolicyStart(w http.ResponseWriter, r *http.Request) {
	html.ExecuteTemplate(w, "start.html", nil)
}

func PostPolicyStart(w http.ResponseWriter, r *http.Request) {
	html.ExecuteTemplate(w, "start.html", struct{ Success bool }{true})
}

// Agents
func AgengLowerDecLanding(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from agents")
}
