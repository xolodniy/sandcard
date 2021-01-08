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
	event      chan eventRaw

	deck        []string
	tablePile   []string
	discardPile []string

	tableLog []log
}

type eventRaw struct {
	senderID int
	body     []byte
}

type log struct {
	Timestamp time.Time   `json:"timestamp"`
	EventType string      `json:"eventType"`
	Extra     interface{} `json:"extra"`
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
		event:       make(chan eventRaw),
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

func parseEvent(body []byte) (string, Event, error) {
	var template struct {
		Type  string `json:"type"`
		Event Event  `json:"event"`
	}
	if err := json.Unmarshal(body, &template); err != nil {
		return "", nil, errors.New("invalid event structure")
	}
	return template.Type, template.Event, nil
}

func (t *Table) logEvent(eventType string, event interface{}) {
	var (
		now = time.Now()
		l   = log{
			Timestamp: now,
			EventType: eventType,
			Extra:     event,
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
			t.event <- eventRaw{
				senderID: player.id,
				body:     msg,
			}
		}
	}()
	return nil
}

func (t *Table) UserIsGone(userID int) {
	i, ok := t.userByID(userID)
	if !ok {
		return
	}
	t.discardPile = append(t.discardPile, t.players[i].cards...)
	t.players = append(t.players[:i], t.players[i+1:]...)
}

const (
	h6 = "heart6"
	d6 = "diamond6"
	c6 = "club6"
	s6 = "spade6"
)
