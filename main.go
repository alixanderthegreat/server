package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type Route struct {
	Path     string
	File     string
	MimeType string
}

var (
	port         = ":8080"
	static_dir   = "./static"
	favicon_dir  = static_dir + "/favicon"
	index_file   = static_dir + "/index.html"
	script_file  = static_dir + "/script.js"
	styles_file  = static_dir + "/styles.css"
	favicon_file = favicon_dir + "/favicon.ico"
)

var (
	routes = []Route{
		{
			Path:     "/",
			File:     index_file,
			MimeType: "text/html; charset=utf-8",
		},
		{
			Path:     "/script.js",
			File:     script_file,
			MimeType: "text/javascript",
		},
		{
			Path:     "/styles.css",
			File:     styles_file,
			MimeType: "text/css",
		},
		{
			Path:     "/favicon.ico",
			File:     favicon_file,
			MimeType: "image/x-icon",
		},
	}
)

// serveFile
func serve(w http.ResponseWriter, r *http.Request, file string, mimeType string) {

	log.Printf("request: %s %s %s", r.Method, r.URL, r.RemoteAddr)

	w.Header().Set("Content-Type", mimeType)

	// Read the file content
	content, err := os.ReadFile(file)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	switch file {
	case index_file:
		data := struct {
			App string
		}{
			App: time.Now().Format("15:04:05"),
		}
		content = inject(content, data)
	}

	w.Write(content)
}

func inject(content []byte, data any) []byte {
	tmpl, err := template.New("index").Parse(string(content))
	if err != nil {
		log.Println("Template parsing error:", err)
		return content
	}

	var inject strings.Builder
	err = tmpl.Execute(&inject, data)
	if err != nil {
		log.Println("Template execution error:", err)
		return content
	}

	return []byte(inject.String())
}

func main() {
	for _, route := range routes {
		http.HandleFunc(
			route.Path,
			func(w http.ResponseWriter, r *http.Request) {
				serve(w, r, route.File, route.MimeType)
			},
		)
	}
	log.Println("Serving at " + port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Println("Server failed to start:", err)
	}
}
