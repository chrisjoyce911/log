package main

import (
	log "github.com/chrisjoyce911/log"
)

type User struct {
	Name string
	Age  int
}

func main() {
	// Straight replacement: same API as std log
	log.Println("hello from my log (console)")

	u := User{Name: "Bob", Age: 28}
	log.Printf("user: %+v", u)
}
