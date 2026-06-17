package evaluator

import (
	"ProjetFilViolet/backend/game/deck"
	"sort"
)

type HandRank int

const (
	HighCard HandRank = iota
	OnePair
	TwoPair
	ThreeOfAKind
	Straight
	Flush
	FullHouse
	FourOfAKind
	StraightFlush
	RoyalFlush
)

func (hr HandRank) String() string {
	switch hr {
	case HighCard:
		return "High Card"
	case OnePair:
		return "One Pair"
	case TwoPair:
		return "Two Pair"
	case ThreeOfAKind:
		return "Three of a Kind"
	case Straight:
		return "Straight"
	case Flush:
		return "Flush"
	case FullHouse:
		return "Full House"
	case FourOfAKind:
		return "Four of a Kind"
	case StraightFlush:
		return "Straight Flush"
	case RoyalFlush:
		return "Royal Flush"
	default:
		return "Unknown"
	}
}

// HandScore représente la force d'une main pour comparaison et départage.
type HandScore struct {
	Rank        HandRank
	TieBreakers []int // Valeurs de cartes triées par importance pour départager
}

// Compare renvoie 1 si hs > other, -1 si hs < other, et 0 si hs == other.
func (hs HandScore) Compare(other HandScore) int {
	if hs.Rank > other.Rank {
		return 1
	}
	if hs.Rank < other.Rank {
		return -1
	}
	for i := 0; i < len(hs.TieBreakers) && i < len(other.TieBreakers); i++ {
		if hs.TieBreakers[i] > other.TieBreakers[i] {
			return 1
		}
		if hs.TieBreakers[i] < other.TieBreakers[i] {
			return -1
		}
	}
	return 0
}

// Evaluate5Cards évalue une main de exactement 5 cartes.
func Evaluate5Cards(cards []deck.Card) HandScore {
	if len(cards) != 5 {
		return HandScore{Rank: HighCard, TieBreakers: []int{}}
	}

	// Trier les cartes par valeur décroissante
	sort.Slice(cards, func(i, j int) bool {
		return cards[i].Rank > cards[j].Rank
	})

	isFlush := true
	for i := 1; i < 5; i++ {
		if cards[i].Suit != cards[0].Suit {
			isFlush = false
			break
		}
	}

	// Détecter les quintes (Straight)
	isStraight := false
	isAceLowStraight := false // 5-4-3-2-A
	
	// Cas normal
	if cards[0].Rank-cards[4].Rank == 4 {
		isStraight = true
		// Vérification qu'il n'y a pas de doublon (ce qui empêcherait une quinte mais est géré par la différence de 4 si trié)
		for i := 0; i < 4; i++ {
			if cards[i].Rank == cards[i+1].Rank {
				isStraight = false
			}
		}
	}
	// Cas As-basse quinte (A-5-4-3-2)
	if !isStraight && cards[0].Rank == deck.Ace && cards[1].Rank == deck.Five && cards[2].Rank == deck.Four && cards[3].Rank == deck.Three && cards[4].Rank == deck.Two {
		isStraight = true
		isAceLowStraight = true
	}

	// Compter les occurrences de chaque valeur
	counts := make(map[deck.Rank]int)
	for _, card := range cards {
		counts[card.Rank]++
	}

	type rankCount struct {
		rank  deck.Rank
		count int
	}
	var rcList []rankCount
	for r, c := range counts {
		rcList = append(rcList, rankCount{rank: r, count: c})
	}
	// Trier d'abord par le nombre d'occurrences décroissant, puis par valeur décroissante
	sort.Slice(rcList, func(i, j int) bool {
		if rcList[i].count != rcList[j].count {
			return rcList[i].count > rcList[j].count
		}
		return rcList[i].rank > rcList[j].rank
	})

	// 1. Quinte Flush (et Royal Flush)
	if isFlush && isStraight {
		if isAceLowStraight {
			return HandScore{Rank: StraightFlush, TieBreakers: []int{int(deck.Five)}}
		}
		if cards[0].Rank == deck.Ace {
			return HandScore{Rank: RoyalFlush, TieBreakers: []int{int(deck.Ace)}}
		}
		return HandScore{Rank: StraightFlush, TieBreakers: []int{int(cards[0].Rank)}}
	}

	// 2. Carré (Four of a Kind)
	if rcList[0].count == 4 {
		return HandScore{
			Rank:        FourOfAKind,
			TieBreakers: []int{int(rcList[0].rank), int(rcList[1].rank)},
		}
	}

	// 3. Full House
	if rcList[0].count == 3 && rcList[1].count == 2 {
		return HandScore{
			Rank:        FullHouse,
			TieBreakers: []int{int(rcList[0].rank), int(rcList[1].rank)},
		}
	}

	// 4. Couleur (Flush)
	if isFlush {
		tb := make([]int, 5)
		for i := 0; i < 5; i++ {
			tb[i] = int(cards[i].Rank)
		}
		return HandScore{Rank: Flush, TieBreakers: tb}
	}

	// 5. Quinte (Straight)
	if isStraight {
		if isAceLowStraight {
			return HandScore{Rank: Straight, TieBreakers: []int{int(deck.Five)}}
		}
		return HandScore{Rank: Straight, TieBreakers: []int{int(cards[0].Rank)}}
	}

	// 6. Brelan (Three of a Kind)
	if rcList[0].count == 3 {
		return HandScore{
			Rank:        ThreeOfAKind,
			TieBreakers: []int{int(rcList[0].rank), int(rcList[1].rank), int(rcList[2].rank)},
		}
	}

	// 7. Double Paire (Two Pair)
	if rcList[0].count == 2 && rcList[1].count == 2 {
		return HandScore{
			Rank:        TwoPair,
			TieBreakers: []int{int(rcList[0].rank), int(rcList[1].rank), int(rcList[2].rank)},
		}
	}

	// 8. Paire (One Pair)
	if rcList[0].count == 2 {
		return HandScore{
			Rank:        OnePair,
			TieBreakers: []int{int(rcList[0].rank), int(rcList[1].rank), int(rcList[2].rank), int(rcList[3].rank)},
		}
	}

	// 9. Carte Haute (High Card)
	tb := make([]int, 5)
	for i := 0; i < 5; i++ {
		tb[i] = int(cards[i].Rank)
	}
	return HandScore{Rank: HighCard, TieBreakers: tb}
}

// Evaluate7Cards évalue la meilleure main possible de 5 cartes parmi 7 cartes.
func Evaluate7Cards(cards []deck.Card) HandScore {
	if len(cards) < 5 {
		return HandScore{Rank: HighCard, TieBreakers: []int{}}
	}

	var bestScore HandScore
	first := true

	// Générer toutes les combinaisons de 5 cartes parmi N (pour N=7, il y a 21 combinaisons)
	n := len(cards)
	var comb []deck.Card

	// Fonction récursive simple de combinaison
	var generate func(start int, depth int)
	generate = func(start int, depth int) {
		if depth == 5 {
			score := Evaluate5Cards(append([]deck.Card(nil), comb...))
			if first || score.Compare(bestScore) > 0 {
				bestScore = score
				first = false
			}
			return
		}

		for i := start; i < n; i++ {
			comb = append(comb, cards[i])
			generate(i+1, depth+1)
			comb = comb[:len(comb)-1]
		}
	}

	generate(0, 0)
	return bestScore
}
