package application

import (
	"math/rand"
	"time"
)

const (
	h2, d2, c2, s2 = "heart2", "diamond2", "club2", "spade2"
	h3, d3, c3, s3 = "heart3", "diamond3", "club3", "spade3"
	h4, d4, c4, s4 = "heart4", "diamond4", "club4", "spade4"
	h5, d5, c5, s5 = "heart5", "diamond5", "club5", "spade5"
	h6, d6, c6, s6 = "heart6", "diamond6", "club6", "spade6"
	h7, d7, c7, s7 = "heart7", "diamond7", "club7", "spade7"
	h8, d8, c8, s8 = "heart8", "diamond8", "club8", "spade8"
	h9, d9, c9, s9 = "heart9", "diamond9", "club9", "spade9"

	h10, d10, c10, s10 = "heart10", "diamond10", "club10", "spade10"
	h11, d11, c11, s11 = "heart11", "diamond11", "club11", "spade11"
	h12, d12, c12, s12 = "heart12", "diamond12", "club12", "spade12"
	h13, d13, c13, s13 = "heart13", "diamond13", "club13", "spade13"
	h14, d14, c14, s14 = "heart14", "diamond14", "club14", "spade14"

	jr, jb = "jokerRed", "jokerBlack"
)

var (
	deck36 = Deck{
		h6, d6, c6, s6,
		h7, d7, c7, s7,
		h8, d8, c8, s8,
		h9, d9, c9, s9,
		h10, d10, c10, s10,
		h11, d11, c11, s11,
		h12, d12, c12, s12,
		h13, d13, c13, s13,
		h14, d14, c14, s14,
	}
	deck52 = append(deck36,
		h2, d2, c2, s2,
		h3, d3, c3, s3,
		h4, d4, c4, s4,
		h5, d5, c5, s5,
	)
	deck54 = append(deck52, jr, jb)
)

type Deck []string

func NewDeck(c int) Deck {
	var deck Deck
	switch c {
	case 36:
		deck = deck36
	case 52:
		deck = deck52
	case 54:
		deck = deck54
	default:
		c = 36
		deck = deck36
	}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(deck), func(i, j int) {
		deck[i], deck[j] = deck[j], deck[i]
	})
	return deck
}
