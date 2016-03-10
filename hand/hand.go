package hand

import (
	"fmt"
	"sort"

	"github.com/aultimus/gosouth/card"
)

// Hand represents a collection of cards
type Hand []*card.Card

var numHoleCards = 2
var numCommCards = 5
var sizeHand = 5

// OUTCOME type represents all the possible
// outcomes of a hand
type OUTCOME int

const (
	// H1Win constant
	H1Win = OUTCOME(iota)
	// H2Win constant
	H2Win = OUTCOME(iota)
	// Draw constant
	Draw = OUTCOME(iota)
)

// RANK type represents the hierarchy of winning hands
type RANK int

const (
	// HighCard constant
	HighCard = RANK(iota)
	// OnePair constant
	OnePair = RANK(iota)
	// TwoPair constant
	TwoPair = RANK(iota)
	// ThreeOfAKind constant
	ThreeOfAKind = RANK(iota)
	// Straight constant
	Straight = RANK(iota)
	// Flush constant
	Flush = RANK(iota)
	// FullHouse constant
	FullHouse = RANK(iota)
	// FourOfAKind constant
	FourOfAKind = RANK(iota)
	// StraightFlush constant
	StraightFlush = RANK(iota)
	// RoyalFlush constant
	RoyalFlush = RANK(iota)
)

// Value encapsulates a showdown hand
// and is used to compare showdown hands.
// Kicker may be default value, not all hands have kickers.
// May need more complex expression of RankValue for two pair matchups
// where both players have the same top pair
// Value contains tie breaking logic.
type Value struct {
	Rank      RANK
	RankValue card.RANK
	Kicker    card.RANK
	Hand      Hand
}

// NewHandValue creates a new Value
func NewHandValue(rank RANK, kicker card.RANK, rankValue card.RANK, hand Hand) *Value {
	return &Value{
		Rank:      rank,
		RankValue: rankValue,
		Kicker:    kicker,
		Hand:      hand,
	}
}

func (v Value) String() string {
	return fmt.Sprintf("rank:%d, rank value:%s kicker:%s, hand:%s",
		v.Rank, v.RankValue, v.Kicker, v.Hand)
}

func (h Hand) Len() int           { return len(h) }
func (h Hand) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h Hand) Less(i, j int) bool { return card.RankIndexes[h[i].Rank] < card.RankIndexes[h[j].Rank] }

// FormHand given the hole cards and community cards returns the best
// Value that can be formed
func FormHand(h Hand) (*Value, error) {
	var v *Value
	if len(h) != numHoleCards+numCommCards {
		return v, fmt.Errorf("Argument to FormHand should be hand of %d cards, not %d cards",
			numHoleCards+numCommCards, len(h))
	}
	// Straight
	// TODO - straight isnt higher than four of a kind
	hasStraight, straightHigh, isFlush := straight(h)
	if hasStraight {
		if isFlush {
			// Straight flush
			straightVal = StraightFlush
			// Royal Flush
			if straightHigh == card.Ace {
				straightVal = RoyalFlush
			}
			return NewHandValue(straightVal, card.Nil, straightHigh, h), nil
		}
	}

	// Four of a kind
	hasFour, kicker, rankValue := xOfAKind(h, 4)
	if hasFour {
		return NewHandValue(FourOfAKind, kicker, rankValue, h), nil
	}
	// Full house

	// Flush
	hasFlush := flush(h)
	if hasFlush {
		return NewHandValue(Flush, card.Nil, card.Nil, h), nil
	}

	if hasStraight {
		return NewHandValue(Straight, card.Nil, straightHigh, h), nil
	}

	// Three of a kind
	hasThree, kicker, rankValue := xOfAKind(h, 3)
	if hasThree {
		return NewHandValue(ThreeOfAKind, kicker, rankValue, h), nil
	}
	return v, nil
}

// Showdown determines whether h1 wins,
// h2 wins or if they draw
func Showdown(h1, h2 Hand) OUTCOME {
	v1, err := FormHand(h1)
	if err != nil {
		panic(err)
	}
	v2, err := FormHand(h2)
	if err != nil {
		panic(err)
	}
	if v1.Rank > v2.Rank {
		return H1Win
	} else if v2.Rank > v1.Rank {
		return H2Win
	}
	// Ordering of hand type checks from here on doesnt matter,
	// Both have the same hand value
	// Handle tie breaking hand comparison
	// Royal Flush
	if ofSameRank(v1, v2, RoyalFlush) {
		panic(fmt.Errorf("both %s and %s are royal flushes!?? Impossible", v1, v2))
	}

	// Straight & Straight Flush
	if ofSameRank(v1, v2, StraightFlush) {
		outcome, err := straightTieBreak(v1, v2)
		if err != nil {
			panic(err)
		}
		return outcome
	}

	// Flush - TODO

	// Four of a kind
	if ofSameRank(v1, v2, FourOfAKind) {
		// TODO
		//outcome, err := FourOfAKindTieBreak(v1, v2)
		//if err != nil {
		//	panic(err)
		//}
	}

	// TODO
	return Draw
}

// ofSameRank returns true if the two given hand
// Values match the rank r
func ofSameRank(v1, v2 *Value, r RANK) bool {
	return v1.Rank == r && v2.Rank == r
}

func numSuited(h Hand) (card.SUIT, int) {
	var m = map[card.SUIT]int{
		card.Clubs:    0,
		card.Diamonds: 0,
		card.Hearts:   0,
		card.Spades:   0,
	}
	for _, c := range h {
		m[c.Suit]++
	}
	var largestSuit card.SUIT
	var largestCount int
	for s := range m {
		if m[s] > largestCount {
			largestSuit = s
			largestCount = m[s]
		}
	}
	return largestSuit, largestCount
}

func rmDupsOfOtherSuits(h Hand, s card.SUIT) Hand {
	dups := make(map[card.RANK]bool)

	favSuit := make(map[card.RANK]bool)
	for _, v := range h {
		if v.Suit == s {
			favSuit[v.Rank] = true
		}
	}

	for _, v := range h {
		if v.Suit != s && favSuit[v.Rank] {
			dups[v.Rank] = true
		}
	}
	var cleaned Hand
	for _, v := range h {
		if !dups[v.Rank] ||
			v.Suit == s && dups[v.Rank] {
			cleaned = append(cleaned, v)
		}
	}
	return cleaned
}

func findStraight(h Hand) (bool, card.RANK, bool) {
	var c int
	var highest card.RANK
	for i := len(h) - 1; i > 0; i-- {
		current := h[i]
		next := h[i-1]
		if current.Connected(next) {
			if c == 0 {
				highest = current.Rank
			}
			c++
			if c == sizeHand-1 {
				_, suitedCount := numSuited(h)
				isFlush := suitedCount == sizeHand
				return true, highest, isFlush
			}
		} else {
			c = 0
		}
	}
	return false, card.RANK(""), false
}

// straight returns a bool, representing whether a straight exists
// and if so, the highest value of a straight in the given hand
func straight(h Hand) (bool, card.RANK, bool) {
	// remove duplicates of rank not of the most populous suit
	s, _ := numSuited(h)
	h = rmDupsOfOtherSuits(h, s)
	if len(h) < sizeHand {
		return false, card.RANK(""), false
	}
	// sort cards into order and check for straight
	sort.Sort(h)
	hasStraight, highest, isFlush := findStraight(h)
	if hasStraight {
		return hasStraight, highest, isFlush
	}
	// Make ace low if exists and check for wheel
	lastI := len(h) - 1
	last := h[lastI]

	// No Ace means no straight
	if last.Rank != card.Ace {
		return false, card.RANK(""), false
	}

	h = append(h[:lastI])
	h = append(Hand{last}, h...)
	hasStraight, highest, isFlush = findStraight(h)
	return hasStraight, highest, isFlush
}

func straightTieBreak(v1, v2 *Value) (OUTCOME, error) {
	if !ofSameRank(v1, v2, StraightFlush) ||
		!ofSameRank(v1, v2, Straight) {
		return Draw, fmt.Errorf("one of %s, %s is not a straight or a straight flush", v1, v2)
	}
	if v1.RankValue > v2.RankValue {
		return H1Win, nil
	} else if v2.RankValue > v1.RankValue {
		return H2Win, nil
	} else {
		return Draw, nil // Is this even possible?
	}
}

func flush(h Hand) bool {
	_, count := numSuited(h)
	return count == sizeHand
}

func xOfAKind(h Hand, x int) (bool, card.RANK, card.RANK) {
	rankMap := make(map[card.RANK]int)
	var rank card.RANK
	for _, v := range h {
		rankMap[v.Rank]++
		if rankMap[v.Rank] > rankMap[rank] ||
			rankMap[v.Rank] == rankMap[rank] &&
				card.RankIndexes[v.Rank] > card.RankIndexes[rank] {
			rank = v.Rank
		}
	}
	// should we handle 4 of a kind case when arg is 3?
	if rankMap[rank] != x {
		return false, card.Nil, card.Nil
	}

	kicker := card.Nil
	for _, c := range h {
		if card.RankIndexes[c.Rank] > card.RankIndexes[kicker] && c.Rank != rank {
			kicker = c.Rank
		}
	}

	return true, kicker, rank
}
