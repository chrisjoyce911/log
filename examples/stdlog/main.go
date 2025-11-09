package main

import (
	"log"
)

type User struct {
	Name string
	Age  int
}

func main() {
	log.Println("hello from std log")

	u := User{Name: "Alice", Age: 30}
	log.Printf("user: %+v", u)
}
