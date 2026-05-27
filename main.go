// Podplane <https://podplane.dev>
// Copyright 2026 Nadrama Pty Ltd
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"bytes"
	"embed"
	"html/template"
	"log"
	"net/http"
	"os"
)

//go:embed index.html index.css podplane-logo.svg favicon.svg favicon.ico
var content embed.FS

func main() {
	addr := ":" + getEnv("PORT", "8080")
	message := getEnv("HELLO_MESSAGE", "Hello, World!")
	indexHTML, err := renderIndex(message)
	if err != nil {
		log.Fatal(err)
	}

	files := http.FileServer(http.FS(content))
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" && r.URL.Path != "/index.html" {
			files.ServeHTTP(w, r)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if _, err := w.Write(indexHTML); err != nil {
			log.Printf("write index: %v", err)
		}
	})
	log.Printf("serving hello page on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}

// renderIndex renders the embedded index page with message as the page heading.
func renderIndex(message string) ([]byte, error) {
	index := template.Must(template.ParseFS(content, "index.html"))
	var page bytes.Buffer
	err := index.Execute(&page, struct {
		Message string
	}{
		Message: message,
	})
	return page.Bytes(), err
}

// getEnv returns the environment variable for key, or fallback when it is unset.
func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
