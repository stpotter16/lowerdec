package main

import (
	"embed"
	"html"
	"net/http"
)

var (
    //go:embed all:templates/*
    templatesFS embed.FS
)

func LowerDecLanding (w http.ResponseWriter, r *http.Request) {
    html.Ex
}
