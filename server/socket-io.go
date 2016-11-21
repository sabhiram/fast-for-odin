package server

import (
	"log"
	"time"

	socketio "github.com/googollee/go-socket.io"

	"github.com/sabhiram/fast-for-odin/room"
)

var (
	counter = 0
	names   = []string{"enney", "miiney", "miney", "moe"}
)

func SocketIOHandler() (*socketio.Server, error) {
	io, err := socketio.NewServer(nil)
	if err != nil {
		return nil, err
	}

	counter = 0

	r, err := room.NewRoom(io, "Awesome room", 3, 10*time.Second, 100*time.Second)

	io.On("connection", func(so socketio.Socket) {
		log.Println("on connection")

		so.Join("1234")

		p := room.NewPlayer(names[counter], 10*time.Second)
		counter += 1
		r.AddPlayer(p)

		so.On("message", func(msg string) {
			log.Println("emit:", so.Emit("message", msg))
			so.BroadcastTo("1234", "message", msg)
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
