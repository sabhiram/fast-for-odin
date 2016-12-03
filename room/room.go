package room

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"strings"
	"time"

	socketio "github.com/googollee/go-socket.io"
)

func GetRandomRoomID() (string, error) {
	rb := make([]byte, 4)
	_, err := rand.Read(rb)
	if err != nil {
		return "", err
	}

	ret := base64.URLEncoding.EncodeToString(rb)
	ret = strings.Replace(ret, "=", "0", -1)
	return ret, nil
}

type Player struct {
	Name     string        `json:"Name"`
	TimeLeft time.Duration `json:"TimeLeft"`
}

func NewPlayer(name string, t time.Duration) *Player {
	return &Player{
		Name:     name,
		TimeLeft: t,
	}
}

type Room struct {
	socket *socketio.Server // Socket to bind to this room

	RoomID   string `json:"RoomID"`   // Unique id of this room
	RoomName string `json:"RoomName"` // Name of this room

	Players       []*Player `json:"Players"`       // List of players in the room
	NumPlayers    int       `json:"NumPlayers"`    // Number of players expected in the room
	CurrentPlayer int       `json:"CurrentPlayer"` // Current player index

	TimeStart    time.Duration `json:"TimeStart"`    // Time each player starts with
	TimePerRound time.Duration `json:"TimePerRound"` // Time added each round

	RoundNum int // Current "round"
}

func NewRoom(socket *socketio.Server, name string, count int, start, perRound time.Duration) (*Room, error) {
	if count <= 0 {
		return nil, errors.New("need 0 or more players")
	}

	id, err := GetRandomRoomID()
	if err != nil {
		return nil, err
	}

	return &Room{
		socket: socket,

		RoomID:   id,
		RoomName: name,

		Players:    []*Player{},
		NumPlayers: count,

		CurrentPlayer: -1, // game not started
		RoundNum:      -1, // game not started

		TimeStart:    start,
		TimePerRound: perRound,
	}, nil
}

func (r *Room) GetID() string {
	return r.RoomID
}

func (r *Room) AddPlayer(p *Player) error {
	existingPlayer := false
	for _, player := range r.Players {
		if player.Name == p.Name {
			existingPlayer = true
		}
	}

	if !existingPlayer {
		r.Players = append(r.Players, p)
	}

	// Update the room for everyone
	return r.BroadcastRoomUpdate()
}

func (r *Room) BroadcastUpdate(msg string) error {
	log.Printf("BroadcastUpdate: %v\n", msg)
	r.socket.BroadcastTo(r.RoomID, "update-room", msg)
	return nil
}

func (r *Room) BroadcastRoomUpdate() error {
	bs, err := json.Marshal(r)
	if err != nil {
		return err
	}
	return r.BroadcastUpdate(string(bs))
}
