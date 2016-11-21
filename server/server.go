package server

import (
	"net/http"

	socketio "github.com/googollee/go-socket.io"
	"github.com/gorilla/mux"
)

type Server struct {
	*http.Server // composed of a http.Server instance

	io      *socketio.Server // ... which contains an instance of a socket.io thingy
	address string           // server host address
}

func New(address string) (*Server, error) {
	router := mux.NewRouter()

	// Application interface
	router.HandleFunc("/api/create/room", HTTPPostNewRoomHandler())

	// Socket.IO
	sio, err := SocketIOHandler()
	if err != nil {
		return nil, err
	}
	router.Handle("/socket.io/", sio)

	// Handle static routes
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static"))) // TODO:

	return &Server{
		Server: &http.Server{
			Addr:    address,
			Handler: router,
		},

		io:      sio,
		address: address,
	}, nil
}
