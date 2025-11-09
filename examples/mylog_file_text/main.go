package main

import (
	log "github.com/chrisjoyce911/log"
)

type User struct {
	Name string
	Age  int
}

func main() {
	f, err := log.SetOutputFile("logs/app.log")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	log.Println("hello from my log (file text)")

	u := User{Name: "Carol", Age: 35}
	log.Printf("user: %+v", u)
}
