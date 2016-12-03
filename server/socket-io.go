package server

import (
	"log"
	// "time"
	"encoding/json"

	socketio "github.com/googollee/go-socket.io"

	// "github.com/sabhiram/fast-for-odin/room"
)

var (
	counter = 0
)

func (s *Server) SocketIOHandler() (*socketio.Server, error) {
	io, err := socketio.NewServer(nil)
	if err != nil {
		return nil, err
	}

	counter = 0

	io.On("connection", func(so socketio.Socket) {
		roomID := ""

		so.On("roomid", func(id string) {
			roomID = id
			room, err := s.GetRoom(id)
			if err != nil {
				log.Printf("Error: %s\n", err.Error())
				return
			}

			so.Join(roomID)
			bs, err := json.Marshal(room)
			if err != nil {
				log.Printf("Error: %s\n", err.Error())
				return
			}
			so.Emit("update-room", string(bs))
		})

		so.On("message", func(msg string) {
			log.Println("emit:", so.Emit("message", msg))
			if len(roomID) > 0 {
				so.BroadcastTo(roomID, "message", msg)
			}
		})

		so.On("disconnection", func() {
			log.Println("on disconnect")
			counter -= 1
		})
	})

	io.On("error", func(so socketio.Socket, err error) {
		log.Println("error:", err)
	})

	return io, nil
}
