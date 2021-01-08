package application

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var (
	tableIDCounter  int
	playerIDCounter int
)

type Table struct {
	ID         int
	maxPlayers int
	players    []Player
	event      chan event

	deck        []string
	tablePile   []string
	discardPile []string

	tableLog []log
}

type event struct {
	senderID int
	body     []byte
}

type log struct {
	timestamp time.Time
	message   string
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
		event:       make(chan event),
		players:     make([]Player, 0),
		deck:        make([]string, 0),
		tablePile:   make([]string, 0),
		discardPile: make([]string, 0),
		tableLog:    make([]log, 0),
	}
}

func (t *Table) AddDeck(c int) *Table {
	t.deck = append(t.deck, h6, d6, c6, s6)
	return t
}

func (t *Table) Start() {
	for ev := range t.event {
		evType, event, err := parseEvent(ev.body)
		if err != nil {
			t.sayTo(ev.senderID, err)
			continue
		}
		t.handleEvent(evType, event, ev.senderID)
	}
}

func (t *Table) sayTo(userID int, message interface{}) {
	i, ok := t.userByID(userID)
	if ok {
		err := t.players[i].connection.WriteJSON(message)
		if err != nil {
			logrus.WithError(err).Error("can't write message to connection")
		}
	}
}

func parseEvent(body []byte) (string, interface{}, error) {
	var template struct {
		Type  string      `json:"type"`
		Event interface{} `json:"event"`
	}
	if err := json.Unmarshal(body, &template); err != nil {
		return "", nil, errors.New("invalid event structure")
	}
	return template.Type, template.Event, nil
}

func (t *Table) sayAllPlayers(message string) {
	var (
		now = time.Now()
		l   = log{
			timestamp: now,
			message:   message,
		}
	)
	t.tableLog = append(t.tableLog, l)
	for i := range t.players {
		if err := t.players[i].connection.WriteJSON(l); err != nil {
			logrus.WithError(err).Error("can't write message to socket")
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
			t.event <- event{
				senderID: player.id,
				body:     msg,
			}
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
