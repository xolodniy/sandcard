package application

import (
	"errors"
	"fmt"
)

const (
	eventTypeGetCardFromDeck = "get_card_from_deck"
	eventTypeShowCards       = "show_cards"
)

func (t *Table) handleEvent(evType string, event interface{}, senderID int) {
	switch evType {
	case eventTypeGetCardFromDeck:
		t.sayTo(senderID, t.GetCardFromTable(senderID))
	case eventTypeShowCards:
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
	return "done"
}

func (t *Table) RetrieveCards(userID int) interface{} {
	i, ok := t.userByID(userID)
	if !ok {
		return errors.New("user not found")
	}
	return t.players[i].cards
}

func (t *Table) userByID(userID int) (int, bool) {
	for i := range t.players {
		if t.players[i].id == userID {
			return i, true
		}
	}
	return 0, false
}
