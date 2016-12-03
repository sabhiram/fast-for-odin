package server

import (
	"net/http"
	"sync"

	socketio "github.com/googollee/go-socket.io"
	"github.com/gorilla/mux"

	"github.com/sabhiram/fast-for-odin/room"
)

type ErrRoomNotFound struct{}

func (e *ErrRoomNotFound) Error() string {
	return "room not found"
}

type ErrRoomExists struct{}

func (e *ErrRoomExists) Error() string {
	return "room already exists"
}

type Server struct {
	*http.Server // composed of a http.Server instance

	io      *socketio.Server // ... which contains an instance of a socket.io thingy
	address string           // server host address

	roomLock sync.RWMutex
	rooms    map[string]*room.Room // map of rooms
}

func New(address string) (*Server, error) {
	router := mux.NewRouter()

	s := &Server{
		Server: &http.Server{
			Addr: address,
		},
		address: address,
		rooms:   make(map[string]*room.Room),
	}

	sio, err := s.SocketIOHandler()
	if err != nil {
		return nil, err
	}
	s.io = sio

	// Application interface
	router.HandleFunc("/api/create/room", s.HTTPPostNewRoomHandler())

	// Handle socket.io stuff
	router.Handle("/socket.io/", sio)

	// Handle static routes
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static"))) // TODO:

	s.Handler = router
	return s, nil
}

func (s *Server) RegisterRoom(r *room.Room) error {
	s.roomLock.Lock()
	defer s.roomLock.Unlock()

	if _, ok := s.rooms[r.GetID()]; ok {
		return &ErrRoomExists{}
	} else {
		s.rooms[r.GetID()] = r
	}
	return nil
}

func (s *Server) GetRoom(id string) (*room.Room, error) {
	s.roomLock.RLock()
	defer s.roomLock.RUnlock()

	if r, ok := s.rooms[id]; ok {
		return r, nil
	}
	return nil, &ErrRoomNotFound{}
}
