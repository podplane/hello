// Podplane <https://podplane.dev>
// Copyright 2026 Nadrama Pty Ltd
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"bytes"
	"embed"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"unicode/utf8"
)

const maxMessageFileSize = 200 * 1024

//go:embed index.html index.css podplane-logo.svg favicon.svg favicon.ico
var content embed.FS

func main() {
	addr := ":" + getEnv("PORT", "8080")
	message := helloMessage()
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

// helloMessage returns HELLO_MESSAGE, or the default message when it is unset.
// When HELLO_MESSAGE starts with '/', it is treated as a file path and the
// rendered message is either the file contents or a safe-to-display error.
func helloMessage() string {
	message := getEnv("HELLO_MESSAGE", "Hello, World!")
	if !strings.HasPrefix(message, "/") {
		return message
	}

	file, err := os.Open(message)
	if err != nil {
		return err.Error()
	}
	defer file.Close()

	contents, err := io.ReadAll(io.LimitReader(file, maxMessageFileSize+1))
	if err != nil {
		return "read " + message + ": " + err.Error()
	}
	if len(contents) > maxMessageFileSize {
		return "read " + message + ": file is larger than 200 KiB"
	}
	if !utf8.Valid(contents) || strings.ContainsFunc(string(contents), func(r rune) bool {
		return r < ' ' && r != '\t' && r != '\n' && r != '\r'
	}) {
		return "read " + message + ": file contains binary or non-UTF-8 data"
	}
	return string(contents)
}
