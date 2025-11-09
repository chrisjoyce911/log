package main

import (
	"io"
	"os"

	log "github.com/chrisjoyce911/log"
)

type User struct {
	Name string
	Age  int
}

func main() {
	f, err := os.Create("app.jsonl")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// Use JSON handler to write JSON lines to the file only.
	log.SetOutput(io.Discard) // disable default text output
	jh := log.NewJSONHandler(f)
	log.AddHandler(log.LevelInfo, jh)

	log.Println("hello from my log (file json)") // goes through handlers

	u := User{Name: "Dora", Age: 41}
	log.Info("user struct", "user", u)
}
