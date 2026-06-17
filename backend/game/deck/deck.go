package deck

import (
	"fmt"
	"math/rand"
	"time"
)

type Suit int

const (
	Spades Suit = iota // ♠
	Hearts             // ♥
	Diamonds           // ♦
	Clubs              // ♣
)

func (s Suit) String() string {
	switch s {
	case Spades:
		return "s"
	case Hearts:
		return "h"
	case Diamonds:
		return "d"
	case Clubs:
		return "c"
	default:
		return "?"
	}
}

type Rank int

const (
	Two Rank = iota + 2
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
	Ace
)

func (r Rank) String() string {
	switch r {
	case Jack:
		return "J"
	case Queen:
		return "Q"
	case King:
		return "K"
	case Ace:
		return "A"
	default:
		return fmt.Sprintf("%d", r)
	}
}

type Card struct {
	Rank Rank `json:"rank"`
	Suit Suit `json:"suit"`
}

func (c Card) String() string {
	return c.Rank.String() + c.Suit.String()
}

type Deck struct {
	Cards []Card
}

func NewDeck() *Deck {
	d := &Deck{
		Cards: make([]Card, 0, 52),
	}
	for suit := Spades; suit <= Clubs; suit++ {
		for rank := Two; rank <= Ace; rank++ {
			d.Cards = append(d.Cards, Card{Rank: rank, Suit: suit})
		}
	}
	return d
}

func (d *Deck) Shuffle() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(d.Cards), func(i, j int) {
		d.Cards[i], d.Cards[j] = d.Cards[j], d.Cards[i]
	})
}

func (d *Deck) Draw() (Card, error) {
	if len(d.Cards) == 0 {
		return Card{}, fmt.Errorf("deck is empty")
	}
	card := d.Cards[0]
	d.Cards = d.Cards[1:]
	return card, nil
}
