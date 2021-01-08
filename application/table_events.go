package application

import (
	"errors"
	"fmt"
)

const (
	eventTypeGetCardFromDeck = "get_card_from_deck"
	eventTypeShowTableInfo   = "show_table_info"
)

type Event map[string]interface{}

func (t *Table) handleEvent(evType string, event Event, senderID int) {
	switch evType {
	case eventTypeGetCardFromDeck:
		t.sayTo(senderID, t.GetCardFromTable(senderID))
	case eventTypeShowTableInfo:
		t.sayTo(senderID, t.RetrieveCards(senderID))
	default:
		t.sayTo(senderID, "unexpected event type")
	}
}

func (t *Table) GetCardFromTable(userID int) interface{} {
	i, ok := t.userByID(userID)
	if !ok {
		return fmt.Errorf("user not found")
	}
	if len(t.deck) == 0 {
		return errors.New("deck on table is empty")
	}
	deckLastID := len(t.deck) - 1
	t.players[i].cards = append(t.players[i].cards, t.deck[deckLastID])
	t.deck = t.deck[:deckLastID]
	t.logEvent(eventTypeGetCardFromDeck, Event{"userID": userID})
	return "done"
}

func (t *Table) RetrieveCards(userID int) interface{} {
	type rPlayer struct {
		ID         int `json:"id"`
		CardsCount int `json:"cardsCount"`
	}
	type response struct {
		Cards         []string  `json:"cards"`
		Players       []rPlayer `json:"players"`
		DeckCardCount int       `json:"deckCardCount"`
		TablePile     []string  `json:"tablePile"`
		DiscardPile   []string  `json:"discardPile"`
		Log           []log     `json:"log"`
	}
	i, ok := t.userByID(userID)
	if !ok {
		return errors.New("user not found")
	}
	players := make([]rPlayer, len(t.players))
	for i := range players {
		players[i] = rPlayer{
			ID:         t.players[i].id,
			CardsCount: len(t.players[i].cards),
		}
	}
	return response{
		Cards:         t.players[i].cards,
		Players:       players,
		DeckCardCount: len(t.deck),
		TablePile:     t.tablePile,
		DiscardPile:   t.discardPile,
		Log:           t.tableLog,
	}
}

func (t *Table) userByID(userID int) (int, bool) {
	for i := range t.players {
		if t.players[i].id == userID {
			return i, true
		}
	}
	return 0, false
}
