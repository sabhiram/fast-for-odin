package main

import (
	"log"

	"github.com/sabhiram/fast-for-odin/server"
)

func main() {
	s, err := server.New("0.0.0.0:5000")
	if err != nil {
		log.Fatalf("Error: %s\n", err.Error())
	}

	// Run the server!
	log.Printf("Server listening at port 5000\n")
	log.Fatal(s.ListenAndServe())
}

func init() {
	log.SetFlags(0)
	log.SetPrefix("")
}
