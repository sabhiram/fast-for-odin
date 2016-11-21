package server

import (
	"log"
	"net/http"
)

func HTTPPostNewRoomHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("PostNewRoomHandler called\n")
		w.Write([]byte{'d', 'o', 'n', 'e'})
	}
}
