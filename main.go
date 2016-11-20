package main

import (
	"log"
	"net/http"

	socketio "github.com/googollee/go-socket.io"
	"github.com/gorilla/mux"
)

// handleErr will issue a fatal on any error.
func handleErr(err error) {
	if err != nil {
		log.Fatalf("Error: %s\n", err.Error())
	}
}

func PostNewRoomHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("PostNewRoomHandler called\n")
		w.Write([]byte{'d', 'o', 'n', 'e'})
	}
}

func SocketIOHandler() *socketio.Server {
	io, err := socketio.NewServer(nil)
	handleErr(err)

	io.On("connection", func(so socketio.Socket) {
		log.Println("on connection")

		so.Join("room0")

		so.On("message", func(msg string) {
			log.Println("emit:", so.Emit("message", msg))
			so.BroadcastTo("room0", "message", msg)
		})

		so.On("disconnection", func() {
			log.Println("on disconnect")
		})
	})

	io.On("error", func(so socketio.Socket, err error) {
		log.Println("error:", err)
	})

	return io
}

func main() {

	// Make a router, register some routes
	router := mux.NewRouter()

	// Application interface
	router.HandleFunc("/api/create/room", PostNewRoomHandler())

	// Socket.IO
	router.Handle("/socket.io/", SocketIOHandler())

	// Handle static routes
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static"))) // TODO:

	// Run the server!
	log.Printf("Server listening at port 5000\n")
	log.Fatal(http.ListenAndServe(":5000", router))
}

func init() {
	log.SetFlags(0)
	log.SetPrefix("")
}
