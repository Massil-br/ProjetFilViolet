package deck

import (
	"testing"
)

func TestNewDeck(t *testing.T) {
	d := NewDeck()
	if len(d.Cards) != 52 {
		t.Errorf("Expected deck size 52, got %d", len(d.Cards))
	}
}

func TestShuffle(t *testing.T) {
	d1 := NewDeck()
	d2 := NewDeck()
	d2.Shuffle()

	// They shouldn't be in the exact same order after shuffle (extremely low probability)
	same := true
	for i := range d1.Cards {
		if d1.Cards[i] != d2.Cards[i] {
			same = false
			break
		}
	}
	if same {
		t.Errorf("Deck order did not change after shuffle")
	}
}

func TestDraw(t *testing.T) {
	d := NewDeck()
	card, err := d.Draw()
	if err != nil {
		t.Errorf("Unexpected error drawing card: %v", err)
	}
	if len(d.Cards) != 51 {
		t.Errorf("Expected remaining cards 51, got %d", len(d.Cards))
	}
	if card.Rank != Two || card.Suit != Spades {
		t.Errorf("Expected first card to be 2 of Spades, got %s", card.String())
	}
}
