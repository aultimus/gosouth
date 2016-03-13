package hand

import (
	"fmt"
	"sort"

	"github.com/aultimus/gosouth/card"
	"github.com/aultimus/gosouth/deck"
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
// Value contains tie breaking logic.
type Value struct {
	Rank RANK
	Hand Hand
}

// NewHandValue creates a new Value
func NewHandValue(rank RANK, hand Hand) *Value {
	return &Value{
		Rank: rank,
		Hand: hand,
	}
}

func tieBreak(v1, v2 *Value) OUTCOME {
	if v1.Rank != v2.Rank {
		panic(fmt.Errorf("%s & %s are not of same rank", v1, v2))
	}
	// assuming hand is sorted correctly by detection function
	for i := range v1.Hand {
		r1 := card.RankIndexes[v1.Hand[i].Rank]
		r2 := card.RankIndexes[v2.Hand[i].Rank]
		if r1 > r2 {
			return H1Win
		} else if r2 > r1 {
			return H2Win
		}
	}
	return Draw
}

func (v Value) String() string {
	return fmt.Sprintf("rank:%d, hand:%s",
		v.Rank, v.Hand)
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
	hasStraight, isFlush, formedHand := straight(h)
	if hasStraight {
		if isFlush {
			// Straight flush
			straightVal := StraightFlush
			// Royal Flush
			if formedHand[sizeHand-1].Rank == card.Ace {
				straightVal = RoyalFlush
			}
			return NewHandValue(straightVal, formedHand), nil
		}
	}

	// Four of a kind
	hasFour, formedHand := xOfAKind(h, 4)
	if hasFour {
		return NewHandValue(FourOfAKind, formedHand), nil
	}

	// Full house
	hasFullHouse, formedHand := fullHouse(h)
	if hasFullHouse {
		return NewHandValue(FullHouse, formedHand), nil
	}

	// Flush
	hasFlush, formedHand := flush(h)
	if hasFlush {
		return NewHandValue(Flush, formedHand), nil
	}

	// Straight
	if hasStraight {
		return NewHandValue(Straight, formedHand), nil
	}

	// Three of a kind
	hasThree, formedHand := xOfAKind(h, 3)
	if hasThree {
		return NewHandValue(ThreeOfAKind, formedHand), nil
	}

	// Two Pair
	hasTwoPair, formedHand := twoPair(h)
	if hasTwoPair {
		return NewHandValue(TwoPair, formedHand), nil
	}

	// Pair
	hasOnePair, formedHand := onePair(h)
	if hasOnePair {
		return NewHandValue(OnePair, formedHand), nil
	}

	formedHand = highCard(h)
	return NewHandValue(HighCard, formedHand), nil
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
	return tieBreak(v1, v2)
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

func findStraight(h Hand) (bool, bool, Hand) {
	var c int
	var formedHand Hand
	for i := len(h) - 1; i > 0; i-- {
		current := h[i]
		next := h[i-1]
		if current.Connected(next) {
			if c == 0 {
				formedHand = Hand{current}
			}
			formedHand = append(Hand{next}, formedHand...)
			c++
			if c == sizeHand-1 {
				_, suitedCount := numSuited(formedHand)
				isFlush := suitedCount == sizeHand
				return true, isFlush, formedHand
			}
		} else {
			c = 0
		}
	}
	return false, false, formedHand
}

// straight returns a bool, representing whether a straight exists
// and if so, the highest value of a straight in the given hand
func straight(h Hand) (bool, bool, Hand) {
	var formedHand Hand
	// remove duplicates of rank not of the most populous suit
	s, _ := numSuited(h)
	h = rmDupsOfOtherSuits(h, s)
	if len(h) < sizeHand {
		return false, false, formedHand
	}
	// sort cards into order and check for straight
	sort.Sort(h)
	hasStraight, isFlush, formedHand := findStraight(h)
	if hasStraight {
		return hasStraight, isFlush, formedHand
	}
	// Make ace low if exists and check for wheel
	lastI := len(h) - 1
	last := h[lastI]

	// No Ace means no straight
	if last.Rank != card.Ace {
		return false, false, formedHand
	}
	h = append(h[:lastI])
	h = append(Hand{last}, h...)
	return findStraight(h)
}

// flush returns true if a flush is present and
//  the five highest cards of that suit
func flush(h Hand) (bool, Hand) {
	// TODO: Sort cards and return sorted hand
	s, count := numSuited(h)
	if count < sizeHand {
		return false, h
	}

	var f Hand
	for _, c := range h {
		if c.Suit == s {
			f = append(f, c)
		}
	}
	// sort hand
	sort.Sort(f)
	// take last 5
	return true, f[len(f)-sizeHand:]
}

// rankFreqMap returns a map of rank to frequency
// for the given hand.
func rankFreqMap(h Hand) map[card.RANK]int {
	rankMap := make(map[card.RANK]int)
	for _, v := range h {
		rankMap[v.Rank]++
	}
	return rankMap
}

// findXOfAKind returns true if x cards of
// the same rank are present in the given Hand
// Returns a hand composed of the given cards and
// the original hand with the given cards removed.
func findXOfAKind(h Hand, x int) (bool, Hand, Hand) {
	rankMap := rankFreqMap(h)
	topRank := card.Nil
	for r, c := range rankMap {
		if c > rankMap[topRank] ||
			c == rankMap[topRank] &&
				card.RankIndexes[r] > card.RankIndexes[topRank] {
			topRank = r
		}
	}

	// should we handle 4 of a kind case when arg is 3?
	if rankMap[topRank] != x {
		return false, h, h
	}
	var f Hand
	var reduced Hand
	for _, c := range h {
		if c.Rank == topRank {
			f = append(f, c) // add to new hand
		} else {
			reduced = append(reduced, c)
		}
	}
	return true, reduced, f
}

func xOfAKind(h Hand, x int) (bool, Hand) {
	hasX, h, f := findXOfAKind(h, x)
	if !hasX {
		return false, f
	}

	// add kickers to new hand
	numKickers := sizeHand - x
	var k *card.Card
	for i := 0; i < numKickers; i++ {
		h, k = popKicker(h)
		f = append(f, k)
	}

	return true, f
}

// TODO: code in fullHouse and twoPair
// is duplicated somewhat
func fullHouse(h Hand) (bool, Hand) {
	hasX, h, f1 := findXOfAKind(h, 3)
	if !hasX {
		return false, h
	}

	hasX, _, f2 := findXOfAKind(h, 2)
	if !hasX {
		return false, h
	}
	return true, append(f1, f2...)
}

func twoPair(h Hand) (bool, Hand) {
	hasX, h, p1 := findXOfAKind(h, 2)
	if !hasX {
		return false, h
	}
	hasX, h, p2 := findXOfAKind(h, 2)
	if !hasX {
		return false, h
	}
	_, k := popKicker(h)

	return true, append(append(p1, p2...), k)
}

func onePair(h Hand) (bool, Hand) {
	hasX, h, f := findXOfAKind(h, 2)
	if !hasX {
		return false, h
	}
	var k *card.Card
	for i := 0; i < sizeHand-2; i++ {
		h, k = popKicker(h)
		f = append(f, k)
	}
	return true, f
}

func highCard(h Hand) Hand {
	var k *card.Card
	var f Hand
	for i := 0; i < sizeHand; i++ {
		h, k = popKicker(h)
		f = append(f, k)
	}
	return f
}

// TODO: Seperate out detector funcs and utility funcs

// popKicker returns the highest ranked
// card in a given hand and returns the hand
// with the card removed.
func popKicker(h Hand) (Hand, *card.Card) {
	var topCard *card.Card
	var topInd int
	for i, c := range h {
		if i == 0 || card.RankIndexes[c.Rank] > card.RankIndexes[topCard.Rank] {
			topCard = c
			topInd = i
		}
	}
	return append(h[:topInd], h[topInd+1:]...), topCard
}

// Remove removes a card from the given hand
func Remove(h Hand, c *card.Card) Hand {
	d := deck.Deck(h)
	d, err := deck.Remove(d, c)
	h = Hand(d)
	if err != nil {
		panic(err)
	}
	return h
}
