package engine

import (
	"ProjetFilViolet/backend/game/deck"
	"ProjetFilViolet/backend/game/evaluator"
	"fmt"
	"sort"
	"time"
)

type Stage int

const (
	Waiting Stage = iota
	PreFlop
	Flop
	Turn
	River
	Showdown
)

func (s Stage) String() string {
	switch s {
	case Waiting:
		return "Waiting"
	case PreFlop:
		return "Pre-Flop"
	case Flop:
		return "Flop"
	case Turn:
		return "Turn"
	case River:
		return "River"
	case Showdown:
		return "Showdown"
	default:
		return "Unknown"
	}
}

type Player struct {
	ID                uint        `json:"id"`
	NickName          string      `json:"nick_name"`
	Chips             uint64      `json:"chips"`
	HoleCards         []deck.Card `json:"hole_cards,omitempty"` // Masqué pour les autres joueurs
	CurrentBet        uint64      `json:"current_bet"`
	TotalContribution uint64      `json:"total_contribution"`
	IsActive          bool        `json:"is_active"`
	IsAllIn           bool        `json:"is_all_in"`
	HasActed          bool        `json:"has_acted"`
}

type Pot struct {
	Amount   uint64 `json:"amount"`
	Eligible []uint `json:"eligible"`
}

type Table struct {
	ID                uint        `json:"id"`
	Players           []*Player   `json:"players"`
	Deck              *deck.Deck  `json:"-"`
	CommunityCards    []deck.Card `json:"community_cards"`
	Stage             Stage       `json:"stage"`
	DealerIdx         int         `json:"dealer_idx"`
	SmallBlindIdx     int         `json:"small_blind_idx"`
	BigBlindIdx       int         `json:"big_blind_idx"`
	CurrentTurnIdx    int         `json:"current_turn_idx"`
	CurrentBet        uint64      `json:"current_bet"`
	MinRaise          uint64      `json:"min_raise"`
	SmallBlindAmount  uint64      `json:"small_blind_amount"`
	BigBlindAmount    uint64      `json:"big_blind_amount"`
	GameInProgress    bool        `json:"game_in_progress"`
	Pots              []Pot       `json:"pots"`
	TurnStartTime     time.Time   `json:"turn_start_time"`
	ActionTimeoutSecs int         `json:"action_timeout_secs"`
	LastActionMessage string      `json:"last_action_message"`
}

func NewTable(id uint, smallBlind, bigBlind uint64) *Table {
	return &Table{
		ID:                id,
		Players:           make([]*Player, 0),
		Deck:              deck.NewDeck(),
		CommunityCards:    make([]deck.Card, 0),
		Stage:             Waiting,
		DealerIdx:         -1,
		SmallBlindIdx:     -1,
		BigBlindIdx:       -1,
		CurrentTurnIdx:    -1,
		SmallBlindAmount:  smallBlind,
		BigBlindAmount:    bigBlind,
		GameInProgress:    false,
		ActionTimeoutSecs: 30, // 30 secondes de réflexion par défaut
	}
}

// AddPlayer ajoute un joueur à la table avec un montant de jetons (Buy-in).
func (t *Table) AddPlayer(id uint, nickName string, chips uint64) error {
	if t.GameInProgress {
		return fmt.Errorf("game is in progress, cannot join now")
	}
	if len(t.Players) >= 7 { // Limite par défaut
		return fmt.Errorf("table is full")
	}
	for _, p := range t.Players {
		if p.ID == id {
			return fmt.Errorf("player already at table")
		}
	}
	t.Players = append(t.Players, &Player{
		ID:       id,
		NickName: nickName,
		Chips:    chips,
		IsActive: false,
	})
	t.LastActionMessage = fmt.Sprintf("%s a rejoint la table avec %d jetons", nickName, chips)
	return nil
}

// RemovePlayer retire un joueur de la table et retourne ses jetons.
func (t *Table) RemovePlayer(id uint) (uint64, error) {
	for i, p := range t.Players {
		if p.ID == id {
			chips := p.Chips + p.TotalContribution + p.CurrentBet // Rembourser les jetons engagés dans la main
			if t.GameInProgress && p.IsActive {
				// Si la partie est en cours, le joueur se couche d'abord
				t.FoldPlayer(p.ID)
			}
			t.Players = append(t.Players[:i], t.Players[i+1:]...)
			t.LastActionMessage = fmt.Sprintf("%s a quitté la table", p.NickName)
			
			// Si moins de 2 joueurs restent, on arrête la partie
			if len(t.Players) < 2 && t.GameInProgress {
				t.EndGameEarly()
			}
			return chips, nil
		}
	}
	return 0, fmt.Errorf("player not found")
}

// StartHand commence une nouvelle main de Poker.
func (t *Table) StartHand() error {
	activeCount := 0
	for _, p := range t.Players {
		if p.Chips > 0 {
			activeCount++
		}
	}
	if activeCount < 2 {
		return fmt.Errorf("not enough players with chips to start")
	}

	t.GameInProgress = true
	t.CommunityCards = make([]deck.Card, 0)
	t.Pots = nil
	t.Deck = deck.NewDeck()
	t.Deck.Shuffle()

	// Réinitialiser les états des joueurs
	for _, p := range t.Players {
		if p.Chips > 0 {
			p.IsActive = true
			p.IsAllIn = false
			p.HasActed = false
			p.CurrentBet = 0
			p.TotalContribution = 0
			p.HoleCards = []deck.Card{}
			
			card1, _ := t.Deck.Draw()
			card2, _ := t.Deck.Draw()
			p.HoleCards = []deck.Card{card1, card2}
		} else {
			p.IsActive = false
			p.HoleCards = nil
		}
	}

	// Déterminer le Dealer
	t.advanceDealer()

	// Poser les Blinds
	t.postBlinds()

	t.Stage = PreFlop
	t.MinRaise = t.BigBlindAmount

	// UTG (Under the gun) parle en premier en Pre-Flop (gauche de Big Blind)
	t.CurrentTurnIdx = t.nextActivePlayerIdx(t.BigBlindIdx)
	t.resetHasActedForNewRound(true) // En PreFlop, SB et BB doivent parler

	t.TurnStartTime = time.Now()
	t.LastActionMessage = "Une nouvelle main commence !"
	return nil
}

func (t *Table) advanceDealer() {
	if t.DealerIdx == -1 {
		t.DealerIdx = 0
	} else {
		t.DealerIdx = (t.DealerIdx + 1) % len(t.Players)
		for !t.Players[t.DealerIdx].IsActive {
			t.DealerIdx = (t.DealerIdx + 1) % len(t.Players)
		}
	}
}

func (t *Table) postBlinds() {
	// Small Blind
	sbIdx := t.nextActivePlayerIdx(t.DealerIdx)
	t.SmallBlindIdx = sbIdx
	sbPlayer := t.Players[sbIdx]
	sbAmount := t.SmallBlindAmount
	if sbPlayer.Chips <= sbAmount {
		sbAmount = sbPlayer.Chips
		sbPlayer.IsAllIn = true
	}
	sbPlayer.Chips -= sbAmount
	sbPlayer.CurrentBet = sbAmount
	sbPlayer.TotalContribution = sbAmount

	// Big Blind
	bbIdx := t.nextActivePlayerIdx(sbIdx)
	t.BigBlindIdx = bbIdx
	bbPlayer := t.Players[bbIdx]
	bbAmount := t.BigBlindAmount
	if bbPlayer.Chips <= bbAmount {
		bbAmount = bbPlayer.Chips
		bbPlayer.IsAllIn = true
	}
	bbPlayer.Chips -= bbAmount
	bbPlayer.CurrentBet = bbAmount
	bbPlayer.TotalContribution = bbAmount

	t.CurrentBet = bbAmount
}

func (t *Table) nextActivePlayerIdx(current int) int {
	n := len(t.Players)
	idx := (current + 1) % n
	for i := 0; i < n; i++ {
		if t.Players[idx].IsActive && !t.Players[idx].IsAllIn {
			return idx
		}
		idx = (idx + 1) % n
	}
	return current
}

func (t *Table) resetHasActedForNewRound(isPreFlop bool) {
	for i, p := range t.Players {
		if p.IsActive && !p.IsAllIn {
			// En préflop, SB et BB ont déjà misé mais n'ont pas encore "agi" (décidé de check/raise/fold)
			if isPreFlop && (i == t.SmallBlindIdx || i == t.BigBlindIdx) {
				p.HasActed = false
			} else {
				p.HasActed = false
			}
		}
	}
}

// PlayAction exécute l'action d'un joueur.
func (t *Table) PlayAction(playerID uint, action string, amount uint64) error {
	if !t.GameInProgress {
		return fmt.Errorf("aucune partie en cours")
	}
	currentPlayer := t.Players[t.CurrentTurnIdx]
	if currentPlayer.ID != playerID {
		return fmt.Errorf("ce n'est pas votre tour de jouer")
	}

	var err error
	switch action {
	case "fold":
		err = t.FoldPlayer(playerID)
	case "check":
		err = t.CheckPlayer(playerID)
	case "call":
		err = t.CallPlayer(playerID)
	case "raise":
		err = t.RaisePlayer(playerID, amount)
	default:
		return fmt.Errorf("unknown action: %s", action)
	}

	if err != nil {
		return err
	}

	t.processNextTurn()
	return nil
}

func (t *Table) FoldPlayer(playerID uint) error {
	p := t.getPlayerByID(playerID)
	if p == nil {
		return fmt.Errorf("player not found")
	}
	p.IsActive = false
	p.HoleCards = nil
	t.LastActionMessage = fmt.Sprintf("%s se couche (Fold)", p.NickName)

	// Vérifier si un seul joueur actif reste
	activeCount := 0
	var lastActive *Player
	for _, pl := range t.Players {
		if pl.IsActive {
			activeCount++
			lastActive = pl
		}
	}

	if activeCount == 1 {
		// Le dernier joueur actif gagne le pot immédiatement
		t.awardPotToSingleWinner(lastActive)
		return nil
	}
	return nil
}

func (t *Table) CheckPlayer(playerID uint) error {
	p := t.getPlayerByID(playerID)
	if p == nil {
		return fmt.Errorf("joueur non trouvé")
	}
	if p.CurrentBet != t.CurrentBet {
		return fmt.Errorf("impossible de checker, vous devez vous aligner sur la mise de %d (Call)", t.CurrentBet)
	}
	p.HasActed = true
	t.LastActionMessage = fmt.Sprintf("%s parole (Check)", p.NickName)
	return nil
}

func (t *Table) CallPlayer(playerID uint) error {
	p := t.getPlayerByID(playerID)
	if p == nil {
		return fmt.Errorf("joueur non trouvé")
	}
	diff := t.CurrentBet - p.CurrentBet
	if diff == 0 {
		return t.CheckPlayer(playerID)
	}

	if p.Chips <= diff {
		// All-in call
		added := p.Chips
		p.CurrentBet += added
		p.TotalContribution += added
		p.Chips = 0
		p.IsAllIn = true
	} else {
		p.Chips -= diff
		p.CurrentBet += diff
		p.TotalContribution += diff
	}
	p.HasActed = true
	t.LastActionMessage = fmt.Sprintf("%s suit (Call) à hauteur de %d", p.NickName, p.CurrentBet)
	return nil
}

func (t *Table) RaisePlayer(playerID uint, targetBet uint64) error {
	p := t.getPlayerByID(playerID)
	if p == nil {
		return fmt.Errorf("joueur non trouvé")
	}

	// La relance doit être au moins de (CurrentBet + MinRaise)
	minRequired := t.CurrentBet + t.MinRaise
	if targetBet < minRequired {
		// Si le joueur met tout son tapis mais que c'est inférieur à la relance minimum requise, on l'autorise (all-in)
		if targetBet != p.CurrentBet+p.Chips {
			return fmt.Errorf("la relance doit être d'au moins %d (ou tapis/all-in)", minRequired)
		}
	}

	diff := targetBet - p.CurrentBet
	if p.Chips < diff {
		return fmt.Errorf("pas assez de jetons pour relancer à %d", targetBet)
	}

	// Calculer la nouvelle relance minimum
	raiseSize := targetBet - t.CurrentBet
	if raiseSize > t.MinRaise {
		t.MinRaise = raiseSize
	}

	p.Chips -= diff
	p.CurrentBet = targetBet
	p.TotalContribution += diff
	p.HasActed = true

	if p.Chips == 0 {
		p.IsAllIn = true
		t.LastActionMessage = fmt.Sprintf("%s relance All-In à %d", p.NickName, targetBet)
	} else {
		t.LastActionMessage = fmt.Sprintf("%s relance à %d", p.NickName, targetBet)
	}

	t.CurrentBet = targetBet

	// Réinitialiser "HasActed" pour les autres joueurs actifs
	for _, pl := range t.Players {
		if pl.IsActive && !pl.IsAllIn && pl.ID != playerID {
			pl.HasActed = false
		}
	}

	return nil
}

func (t *Table) processNextTurn() {
	if !t.GameInProgress {
		return
	}

	// Si le tour de mise est complété
	if t.isBettingRoundComplete() {
		t.advanceRound()
		return
	}

	// Passer au joueur actif suivant
	t.CurrentTurnIdx = t.nextActivePlayerIdx(t.CurrentTurnIdx)
	t.TurnStartTime = time.Now()
}

func (t *Table) isBettingRoundComplete() bool {
	activeCount := 0
	matchingCount := 0
	for _, p := range t.Players {
		if !p.IsActive || p.IsAllIn {
			continue
		}
		activeCount++
		if p.HasActed && p.CurrentBet == t.CurrentBet {
			matchingCount++
		}
	}
	if activeCount <= 1 {
		return true
	}
	return matchingCount == activeCount
}

func (t *Table) advanceRound() {
	// Collecter les jetons misés au milieu de la table
	for _, p := range t.Players {
		p.CurrentBet = 0
	}
	t.CurrentBet = 0
	t.MinRaise = t.BigBlindAmount

	// Réinitialiser les actions
	t.resetHasActedForNewRound(false)

	// Compter le nombre de joueurs actifs qui ne sont pas all-in
	playersAbleToAct := 0
	for _, p := range t.Players {
		if p.IsActive && !p.IsAllIn {
			playersAbleToAct++
		}
	}

	// Si moins de 2 joueurs peuvent encore miser (tous les autres sont all-in ou couchés), 
	// on distribue les cartes restantes directement et on va au Showdown
	if playersAbleToAct < 2 {
		t.dealRemainingCommunityCards()
		t.resolveShowdown()
		return
	}

	// Changer d'étape
	switch t.Stage {
	case PreFlop:
		t.Stage = Flop
		// Brûler une carte et distribuer le Flop (3 cartes)
		t.Deck.Draw() // Burn card
		c1, _ := t.Deck.Draw()
		c2, _ := t.Deck.Draw()
		c3, _ := t.Deck.Draw()
		t.CommunityCards = append(t.CommunityCards, c1, c2, c3)
		t.LastActionMessage = "Flop distribué !"
	case Flop:
		t.Stage = Turn
		t.Deck.Draw() // Burn
		c, _ := t.Deck.Draw()
		t.CommunityCards = append(t.CommunityCards, c)
		t.LastActionMessage = "Turn distribué !"
	case Turn:
		t.Stage = River
		t.Deck.Draw() // Burn
		c, _ := t.Deck.Draw()
		t.CommunityCards = append(t.CommunityCards, c)
		t.LastActionMessage = "River distribuée !"
	case River:
		t.resolveShowdown()
		return
	}

	// Après le flop/turn/river, le premier joueur actif à gauche du dealer parle
	t.CurrentTurnIdx = t.nextActivePlayerIdx(t.DealerIdx)
	t.TurnStartTime = time.Now()
}

func (t *Table) dealRemainingCommunityCards() {
	needed := 5 - len(t.CommunityCards)
	for i := 0; i < needed; i++ {
		t.Deck.Draw() // Burn card
		c, err := t.Deck.Draw()
		if err == nil {
			t.CommunityCards = append(t.CommunityCards, c)
		}
	}
}

func (t *Table) resolveShowdown() {
	t.Stage = Showdown
	t.GameInProgress = false

	// Calculer les pots principaux et secondaires
	pots := CalculatePots(t.Players)
	t.Pots = pots

	winnersLog := ""

	// Résoudre chaque pot séparément
	for i, pot := range pots {
		var winners []*Player
		var bestScore evaluator.HandScore
		first := true

		// Évaluer la main de chaque joueur éligible
		for _, pID := range pot.Eligible {
			p := t.getPlayerByID(pID)
			if p == nil || !p.IsActive {
				continue
			}
			score := evaluator.Evaluate7Cards(append(p.HoleCards, t.CommunityCards...))
			if first {
				bestScore = score
				winners = []*Player{p}
				first = false
			} else {
				cmp := score.Compare(bestScore)
				if cmp > 0 {
					bestScore = score
					winners = []*Player{p}
				} else if cmp == 0 {
					winners = append(winners, p)
				}
			}
		}

		if len(winners) > 0 {
			// Partager le pot entre les gagnants
			share := pot.Amount / uint64(len(winners))
			remainder := pot.Amount % uint64(len(winners))

			for idx, w := range winners {
				w.Chips += share
				if idx == 0 {
					w.Chips += remainder // Donne les jetons indivisibles au premier gagnant
				}
				winnersLog += fmt.Sprintf(" %s remporte %d jetons avec %s.", w.NickName, share, bestScore.Rank.String())
			}
		}
		_ = i
	}

	t.LastActionMessage = "Abattage (Showdown) ! " + winnersLog

	// Gérer l'élimination des joueurs sans jetons
	for _, p := range t.Players {
		p.HoleCards = nil // Les cartes sont révélées/jetées
	}
}

func (t *Table) awardPotToSingleWinner(winner *Player) {
	t.Stage = Showdown
	t.GameInProgress = false

	totalPot := uint64(0)
	for _, p := range t.Players {
		totalPot += p.TotalContribution + p.CurrentBet
		p.TotalContribution = 0
		p.CurrentBet = 0
	}

	winner.Chips += totalPot
	t.LastActionMessage = fmt.Sprintf("%s gagne %d jetons suite à l'abandon des adversaires", winner.NickName, totalPot)

	for _, p := range t.Players {
		p.HoleCards = nil
	}
}

func (t *Table) EndGameEarly() {
	t.GameInProgress = false
	t.Stage = Waiting
	// Rembourser les contributions en cours
	for _, p := range t.Players {
		p.Chips += p.TotalContribution + p.CurrentBet
		p.TotalContribution = 0
		p.CurrentBet = 0
		p.HoleCards = nil
	}
}

// CheckTimeouts force le fold si le joueur a dépassé son temps de réflexion
func (t *Table) CheckTimeouts() bool {
	if !t.GameInProgress || t.CurrentTurnIdx == -1 {
		return false
	}
	if time.Since(t.TurnStartTime).Seconds() > float64(t.ActionTimeoutSecs) {
		currentPlayer := t.Players[t.CurrentTurnIdx]
		// Si le check est possible, on check automatiquement, sinon on fold.
		var err error
		if currentPlayer.CurrentBet == t.CurrentBet {
			err = t.CheckPlayer(currentPlayer.ID)
		} else {
			err = t.FoldPlayer(currentPlayer.ID)
		}
		if err == nil {
			t.processNextTurn()
			return true
		}
	}
	return false
}

func (t *Table) getPlayerByID(id uint) *Player {
	for _, p := range t.Players {
		if p.ID == id {
			return p
		}
	}
	return nil
}

// CalculatePots implémente le partage du pot et des pots secondaires
func CalculatePots(players []*Player) []Pot {
	contribs := make(map[uint]uint64)
	for _, p := range players {
		if p.TotalContribution > 0 {
			contribs[p.ID] = p.TotalContribution
		}
	}

	var levels []uint64
	seenLevels := make(map[uint64]bool)
	for _, p := range players {
		if p.TotalContribution > 0 && (p.IsActive || p.IsAllIn) {
			if !seenLevels[p.TotalContribution] {
				seenLevels[p.TotalContribution] = true
				levels = append(levels, p.TotalContribution)
			}
		}
	}
	sort.Slice(levels, func(i, j int) bool { return levels[i] < levels[j] })

	var pots []Pot
	var prevLevel uint64 = 0

	for _, lv := range levels {
		capVal := lv - prevLevel
		potAmount := uint64(0)
		var eligible []uint

		for _, p := range players {
			contrib := contribs[p.ID]
			if contrib > 0 {
				share := contrib
				if share > capVal {
					share = capVal
				}
				potAmount += share
				contribs[p.ID] -= share

				if (p.IsActive || p.IsAllIn) && contrib >= capVal {
					eligible = append(eligible, p.ID)
				}
			}
		}

		if potAmount > 0 && len(eligible) > 0 {
			if len(eligible) == 1 {
				// C'est un pari non suivi (uncalled bet). On le rembourse immédiatement au joueur concerné.
				singlePlayerID := eligible[0]
				for _, p := range players {
					if p.ID == singlePlayerID {
						p.Chips += potAmount
						p.TotalContribution -= potAmount
					}
				}
			} else {
				pots = append(pots, Pot{
					Amount:   potAmount,
					Eligible: eligible,
				})
			}
		}
		prevLevel = lv
	}

	// Rembourser les mises non suivies (uncalled bets)
	for _, p := range players {
		rem := contribs[p.ID]
		if rem > 0 {
			p.Chips += rem
			p.TotalContribution -= rem
		}
	}

	return pots
}
