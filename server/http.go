package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/sabhiram/fast-for-odin/room"
)

type postNewRoomDTO struct {
	UserName   string `json:"UserName"`          // UserName of the creator
	UserID     string `json:"UserID"`            // UserID of the creator
	NumPlayers int    `json:"NumPlayers,string"` // Number of guests
	RoomName   string `json:"RoomName"`          // Name of the room to be (not unique)
}

func (s *Server) HTTPPostNewRoomHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("PostNewRoomHandler called\n")

		decoder := json.NewDecoder(r.Body)

		var pnr postNewRoomDTO
		err := decoder.Decode(&pnr)
		if err != nil {
			log.Printf("Error: %s\n", err.Error())
			return
		}

		defer r.Body.Close()
		log.Printf("%#v\n", pnr)

		nr, err := room.NewRoom(s.io, pnr.RoomName, pnr.NumPlayers, 30*time.Second, 30*time.Second)
		if err != nil {
			log.Printf("Error: %s\n", err.Error())
			return
		}

		err = s.RegisterRoom(nr)

		rsp, err := json.Marshal(nr)
		if err != nil {
			log.Printf("Error: %s\n", err.Error())
			return
		}
		fmt.Fprintf(w, string(rsp))
	}
}
