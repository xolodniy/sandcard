## Documentation about events you can initiate after connect to table

Some events have additional params, which can be passed throw field 'extra' in request
All events expected as JSON struct with fields "type" & "extra"(optional)

### Get card from table deck

Catch top card from the deck. 
All other players will be notified about this event.

 - type: get_card_from_deck
 - extra: none

### Get current table info

Request for the current situation on the table

 - type: show_table_info 
 - extra: none

Example of server response
```bigquery
{
  "cards": ["heard6"], 
  "players": [ 
    {
      "id": 1,
      "cardsCount": 1
    }
  ],
  "deckCardCount": 35,
  "tablePile": [],
  "discardPile": [],
  "log": [
    {
      "timestamp": "2021-01-08T23:02:36.977561+07:00",
      "eventType": "get_card_from_deck",
      "extra": {
        "userID": 1
      }
    }
  ]
}
```

### Reset table 

Drop all history, except logs
Init new card deck on the table (54, 52 or 36 by default cards)
 - type: reset_table
 - extra: {
   cardsCount: 54 int
 }