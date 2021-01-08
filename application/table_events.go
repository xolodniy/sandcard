package application

const (
	eventTypeGetCardFromDeck = "get_card_from_deck"
	eventTypeShowTableInfo   = "show_table_info"
	eventTypeResetTable      = "reset_table"
)

var (
	events = []string{
		eventTypeGetCardFromDeck,
		eventTypeShowTableInfo,
		eventTypeResetTable,
	}
)

type Extra map[string]interface{}

func (t *Table) handleEvent(evType string, event Extra, senderID int) {
	switch evType {
	case eventTypeGetCardFromDeck:
		t.getCardFromDeck(senderID)
	case eventTypeShowTableInfo:
		t.tableInfo(senderID)
	case eventTypeResetTable:
		t.resetTable(senderID, event)
	default:
		t.sayTo(senderID, Extra{
			"error": "unexpected event type",
			"extra": Extra{
				"availableTypes": events,
			},
		})
	}
}

func (t *Table) getCardFromDeck(userID int) {
	if len(t.deck) == 0 {
		t.sayTo(userID, Extra{"error": "deck on table is empty"})
		return
	}
	i, ok := t.userByID(userID)
	if !ok {
		return
	}
	deckLastID := len(t.deck) - 1
	t.players[i].cards = append(t.players[i].cards, t.deck[deckLastID])
	t.deck = t.deck[:deckLastID]
	t.logEvent(eventTypeGetCardFromDeck, Extra{"userID": userID})
}

func (t *Table) tableInfo(userID int) {
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
	players := make([]rPlayer, len(t.players))
	for i := range players {
		players[i] = rPlayer{
			ID:         t.players[i].id,
			CardsCount: len(t.players[i].cards),
		}
	}
	i, ok := t.userByID(userID)
	if !ok {
		return
	}
	t.sayTo(userID, response{
		Cards:         t.players[i].cards,
		Players:       players,
		DeckCardCount: len(t.deck),
		TablePile:     t.tablePile,
		DiscardPile:   t.discardPile,
		Log:           t.tableLog,
	})
}

func (t *Table) resetTable(userID int, event Extra) {
	var cardsCount = 36
	if v, ok := event["cardsCount"]; ok {
		if i, ok := v.(float64); ok {
			cardsCount = int(i)
		}
	}
	t.discardPile = Deck{}
	t.tablePile = Deck{}
	for i := range t.players {
		t.players[i].cards = Deck{}
	}
	t.deck = NewDeck(cardsCount)
	t.logEvent(eventTypeResetTable, Extra{"userID": userID, "cardsCount": cardsCount})
}
