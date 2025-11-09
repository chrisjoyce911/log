package main

import (
	"fmt"
	"io"
	"net/http"

	log "github.com/chrisjoyce911/log"
)

// Logging is HTTP middleware that logs the request and access record using this log package.
// Use package helper for colored method/path highlighting.

func hello(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintln(w, "method not allowed")
		return
	}
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "world"
	}
	fmt.Fprintf(w, "hello %s\n", name)
}

// echo reads the request body and writes it back. Supports POST, PUT, PATCH.
func echo(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost, http.MethodPut, http.MethodPatch:
		b, _ := io.ReadAll(r.Body)
		defer r.Body.Close()
		if len(b) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		_, _ = w.Write(b)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintln(w, "method not allowed")
	}
}

func main() {
	// Colorize console output (auto enables on TTYs; use ColorOn to force).
	log.SetColoredOutput(log.LevelDebug, log.ColorOptions{Mode: log.ColorAuto, ColorLevel: true})
	log.SetFlags(log.Ldate | log.Ltime)

	mux := http.NewServeMux()
	mux.HandleFunc("/hello", hello)
	mux.HandleFunc("/echo", echo)

	// Wrap with logging middleware from the package
	h := log.HTTPLogging(mux, &log.HTTPLogOptions{Mode: log.ColorAuto, IncludeQuery: true, LogPostBody: true, MaxBodyBytes: 2048})

	log.Info("starting server", "addr", ":8080", "routes", "/hello (GET), /echo (POST|PUT|PATCH)")
	if err := http.ListenAndServe(":8080", h); err != nil {
		log.Fatal("server error: ", err)
	}
}
