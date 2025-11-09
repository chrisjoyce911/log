package main

import (
	"io"
	stdlog "log"
	"os"

	log "github.com/chrisjoyce911/log"
)

type User struct {
	Name string
	Age  int
}

func main() {
	// 1) std library logger
	stdlog.Println("[stdlog] hello")
	stdU := User{Name: "Alice", Age: 30}
	stdlog.Printf("[stdlog] user: %+v", stdU)

	// 2) my logger (console)
	log.Println("[mylog-console] hello")
	u1 := User{Name: "Bob", Age: 28}
	log.Printf("[mylog-console] user: %+v", u1)

	// 3) my logger to file (text)
	f, err := os.Create("all_app.log")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	log.SetOutput(f)
	log.Println("[mylog-file-text] hello")
	u2 := User{Name: "Carol", Age: 35}
	log.Printf("[mylog-file-text] user: %+v", u2)

	// 4) my logger to file (JSON)
	jf, err := os.Create("all_app.jsonl")
	if err != nil {
		panic(err)
	}
	defer jf.Close()
	log.SetOutput(io.Discard) // disable default text output
	jh := log.NewJSONHandler(jf)
	log.AddHandler(log.LevelInfo, jh)
	log.Println("[mylog-file-json] hello")
	u3 := User{Name: "Dora", Age: 41}
	log.Info("[mylog-file-json] user struct", "user", u3)
}
