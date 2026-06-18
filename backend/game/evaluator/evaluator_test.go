package evaluator

import (
	"ProjetFilViolet/backend/game/deck"
	"testing"
)

func TestEvaluate5Cards(t *testing.T) {
	// 1. Royal Flush
	rf := []deck.Card{
		{Rank: deck.Ace, Suit: deck.Spades},
		{Rank: deck.King, Suit: deck.Spades},
		{Rank: deck.Queen, Suit: deck.Spades},
		{Rank: deck.Jack, Suit: deck.Spades},
		{Rank: deck.Ten, Suit: deck.Spades},
	}
	score := Evaluate5Cards(rf)
	if score.Rank != RoyalFlush {
		t.Errorf("Expected RoyalFlush, got %s", score.Rank.String())
	}

	// 2. Full House
	fh := []deck.Card{
		{Rank: deck.Ace, Suit: deck.Spades},
		{Rank: deck.Ace, Suit: deck.Hearts},
		{Rank: deck.Ace, Suit: deck.Diamonds},
		{Rank: deck.King, Suit: deck.Clubs},
		{Rank: deck.King, Suit: deck.Spades},
	}
	score = Evaluate5Cards(fh)
	if score.Rank != FullHouse {
		t.Errorf("Expected FullHouse, got %s", score.Rank.String())
	}
	if score.TieBreakers[0] != int(deck.Ace) || score.TieBreakers[1] != int(deck.King) {
		t.Errorf("Expected tiebreakers [14, 13], got %v", score.TieBreakers)
	}

	// 3. Flush
	fl := []deck.Card{
		{Rank: deck.Ace, Suit: deck.Hearts},
		{Rank: deck.Ten, Suit: deck.Hearts},
		{Rank: deck.Eight, Suit: deck.Hearts},
		{Rank: deck.Four, Suit: deck.Hearts},
		{Rank: deck.Two, Suit: deck.Hearts},
	}
	score = Evaluate5Cards(fl)
	if score.Rank != Flush {
		t.Errorf("Expected Flush, got %s", score.Rank.String())
	}

	// 4. Straight (Wheel A-5-4-3-2)
	wheel := []deck.Card{
		{Rank: deck.Ace, Suit: deck.Spades},
		{Rank: deck.Five, Suit: deck.Hearts},
		{Rank: deck.Four, Suit: deck.Diamonds},
		{Rank: deck.Three, Suit: deck.Clubs},
		{Rank: deck.Two, Suit: deck.Spades},
	}
	score = Evaluate5Cards(wheel)
	if score.Rank != Straight {
		t.Errorf("Expected Straight, got %s", score.Rank.String())
	}
	if score.TieBreakers[0] != 5 {
		t.Errorf("Expected high card of Wheel to be 5, got %d", score.TieBreakers[0])
	}
}

func TestEvaluate7Cards(t *testing.T) {
	// 7 Cards containing a Full House and a pair
	cards := []deck.Card{
		{Rank: deck.Ace, Suit: deck.Spades},
		{Rank: deck.Ace, Suit: deck.Hearts},
		{Rank: deck.King, Suit: deck.Spades},
		{Rank: deck.King, Suit: deck.Hearts},
		{Rank: deck.King, Suit: deck.Diamonds},
		{Rank: deck.Two, Suit: deck.Clubs},
		{Rank: deck.Three, Suit: deck.Spades},
	}
	// Best hand should be Full House: Kings full of Aces (K-K-K-A-A)
	score := Evaluate7Cards(cards)
	if score.Rank != FullHouse {
		t.Errorf("Expected best hand to be FullHouse, got %s", score.Rank.String())
	}
	if score.TieBreakers[0] != int(deck.King) || score.TieBreakers[1] != int(deck.Ace) {
		t.Errorf("Expected K-K-K-A-A (13 full of 14), got %v", score.TieBreakers)
	}
}
