package application

import (
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"

	"errors"
	"time"
)

var (
	tableIDCounter  int
	playerIDCounter int
)

type Table struct {
	ID         int
	maxPlayers int
	players    []Player
	event      chan []byte

	deck        []string
	discardPile []string
}

type Player struct {
	id         int
	cards      []string
	connection *websocket.Conn
}

func NewTable() *Table {
	tableIDCounter++
	return &Table{
		ID:          tableIDCounter,
		maxPlayers:  4,
		event:       make(chan []byte),
		players:     make([]Player, 0),
		deck:        make([]string, 0),
		discardPile: make([]string, 0),
	}
}

func (t *Table) AddDeck(c int) *Table {
	t.deck = append(t.deck, h6, d6, c6, s6)
	return t
}

func (t *Table) Start() {
	for ev := range t.event {
		for i := range t.players {
			v := map[string]interface{}{
				"time":    time.Now(),
				"message": string(ev),
			}
			if err := t.players[i].connection.WriteJSON(v); err != nil {
				logrus.WithError(err).Error("can't write message to socket")
			}
		}
	}
}

func (t *Table) Join(c *websocket.Conn) error {
	if len(t.players) == t.maxPlayers {
		return errors.New("max players limit")
	}
	playerIDCounter++
	player := Player{
		id:         playerIDCounter,
		cards:      make([]string, 0),
		connection: c,
	}
	t.players = append(t.players, player)
	go func() {
		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err,
					websocket.CloseGoingAway,
					websocket.CloseNormalClosure,
					websocket.CloseAbnormalClosure,
				) {
					logrus.WithError(err).Error("error read from connection")
				}
				logrus.Debug("lost connection with user ", player.id)
				t.UserIsGone(player.id)
				return
			}
			if len(msg) == 0 {
				logrus.Debug("got empty message from user ", player.id)
				continue
			}
			t.event <- msg
		}
	}()
	return nil
}

func (t *Table) UserIsGone(userID int) {
	for i := range t.players {
		if t.players[i].id == userID {
			t.players = append(t.players[:i], t.players[i+1:]...)
			return
		}
	}
}

const (
	discard = -1
	table   = 0

	h6 = "heart6"
	d6 = "diamond6"
	c6 = "club6"
	s6 = "spade6"
)
