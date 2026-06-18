package engine

import (
	"testing"
)

func TestTableLobby(t *testing.T) {
	table := NewTable(1, 10, 20)
	if len(table.Players) != 0 {
		t.Errorf("Expected 0 players, got %d", len(table.Players))
	}

	err := table.AddPlayer(1, "Alice", 1000)
	if err != nil {
		t.Fatalf("Could not add player Alice: %v", err)
	}

	err = table.AddPlayer(2, "Bob", 1000)
	if err != nil {
		t.Fatalf("Could not add player Bob: %v", err)
	}

	if len(table.Players) != 2 {
		t.Errorf("Expected 2 players, got %d", len(table.Players))
	}

	// Double add
	err = table.AddPlayer(1, "Alice2", 500)
	if err == nil {
		t.Errorf("Expected error adding duplicate player ID, got nil")
	}
}

func TestSimpleHandFlow(t *testing.T) {
	table := NewTable(1, 10, 20)
	_ = table.AddPlayer(1, "Alice", 1000)
	_ = table.AddPlayer(2, "Bob", 1000)

	err := table.StartHand()
	if err != nil {
		t.Fatalf("Failed to start hand: %v", err)
	}

	if table.Stage != PreFlop {
		t.Errorf("Expected PreFlop stage, got %s", table.Stage.String())
	}

	// Small Blind should be Alice (next to Dealer, which defaults to player index 0, so dealer is 0, SB is Bob, wait let's check:
	// DealerIdx = 0 (Alice). SB is next: index 1 (Bob). BB is next: index 0 (Alice).
	// In 2-player game (heads up), SB is usually Dealer, but let's check index:
	// DealerIdx = 0 (Alice).
	// SB is index 1 (Bob), BB is index 0 (Alice).
	// CurrentTurnIdx: UTG is next to BB (Alice, index 0). Next active player to BB is index 1 (Bob).
	// Let's print out values to see who is active and whose turn it is.
	t.Logf("Dealer: %d, SB: %d, BB: %d, Turn: %d", table.DealerIdx, table.SmallBlindIdx, table.BigBlindIdx, table.CurrentTurnIdx)

	currentPlayer := table.Players[table.CurrentTurnIdx]
	t.Logf("Current player turn: %s", currentPlayer.NickName)

	// currentPlayer should Call
	err = table.PlayAction(currentPlayer.ID, "call", 0)
	if err != nil {
		t.Fatalf("Call failed: %v", err)
	}

	// Next player (Alice) should check (she is BB, matched SB's call)
	currentPlayer = table.Players[table.CurrentTurnIdx]
	err = table.PlayAction(currentPlayer.ID, "check", 0)
	if err != nil {
		t.Fatalf("Check failed: %v", err)
	}

	// Stage should advance to Flop
	if table.Stage != Flop {
		t.Errorf("Expected stage to advance to Flop, got %s", table.Stage.String())
	}
}

func TestCalculatePots(t *testing.T) {
	// Let's test a 3-way all-in to verify side pots calculation
	p1 := &Player{ID: 1, NickName: "Alice", Chips: 0, TotalContribution: 100, IsActive: true, IsAllIn: true}
	p2 := &Player{ID: 2, NickName: "Bob", Chips: 0, TotalContribution: 200, IsActive: true, IsAllIn: true}
	p3 := &Player{ID: 3, NickName: "Charlie", Chips: 0, TotalContribution: 300, IsActive: true, IsAllIn: true}

	pots := CalculatePots([]*Player{p1, p2, p3})

	// Pot 1: contribution up to 100 for each. Amount = 100 + 100 + 100 = 300. Eligible: Alice (1), Bob (2), Charlie (3)
	// Pot 2: contribution up to 200 (200-100 = 100 cap). Amount = 0 (Alice has 0 rem) + 100 (Bob) + 100 (Charlie) = 200. Eligible: Bob (2), Charlie (3)
	// Pot 3: Charlie has 100 uncalled bet. It should be refunded. Let's see:
	// Charlie's contribution level was 300. The highest level of other players is 200.
	// So Charlie's 100 uncalled bet is refunded.
	// Let's check!
	if len(pots) != 2 {
		t.Fatalf("Expected 2 pots, got %d", len(pots))
	}

	if pots[0].Amount != 300 {
		t.Errorf("Expected pot 0 amount 300, got %d", pots[0].Amount)
	}
	if len(pots[0].Eligible) != 3 {
		t.Errorf("Expected 3 eligible for pot 0, got %d", len(pots[0].Eligible))
	}

	if pots[1].Amount != 200 {
		t.Errorf("Expected pot 1 amount 200, got %d", pots[1].Amount)
	}
	if len(pots[1].Eligible) != 2 {
		t.Errorf("Expected 2 eligible for pot 1, got %d", len(pots[1].Eligible))
	}

	if p3.Chips != 100 {
		t.Errorf("Expected Charlie to get 100 chips refunded, got %d", p3.Chips)
	}
}
