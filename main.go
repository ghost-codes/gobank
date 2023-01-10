package main

import (
	"fmt"
	"log"
)

func main() {
	fmt.Println("Yeah Buddy!")
	store, err := NewPostgresStore()

	if err != nil {
		log.Fatal(err)
	}

	if err := store.init(); err != nil {
		log.Fatal(err)
	}
	server := NewApiServer(":3000", store)
	server.Run()

}
