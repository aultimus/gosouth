package card

import "fmt"

// Go does not have variant types :(
const (
	NumSuit  = 4
	NumRanks = 13
	NumCards = 4 * 13
)

// RANK represents a card rank "2" -> "A"
type RANK string

const (
	// Nil is used to check for uninitialised ranks
	Nil = RANK("")
	// Two Rank constant
	Two = RANK("2")
	// Three Rank constant
	Three = RANK("3")
	// Four Rank constant
	Four = RANK("4")
	// Five Rank constant
	Five = RANK("5")
	// Six Rank constant
	Six = RANK("6")
	// Seven Rank constant
	Seven = RANK("7")
	// Eight Rank constant
	Eight = RANK("8")
	// Nine Rank constant
	Nine = RANK("9")
	// Ten Rank constant
	Ten = RANK("T")
	// Jack Rank constant
	Jack = RANK("J")
	// Queen Rank constant
	Queen = RANK("Q")
	// King Rank constant
	King = RANK("K")
	// Ace Rank constant
	Ace = RANK("A")
)

// Ranks is a slice of all ranks
var Ranks = []RANK{
	Two,
	Three,
	Four,
	Five,
	Six,
	Seven,
	Eight,
	Nine,
	Ten,
	Jack,
	Queen,
	King,
	Ace,
}

// RankIndexes ...
var RankIndexes = map[RANK]int{
	Two:   0,
	Three: 1,
	Four:  2,
	Five:  3,
	Six:   4,
	Seven: 5,
	Eight: 6,
	Nine:  7,
	Ten:   8,
	Jack:  9,
	Queen: 10,
	King:  11,
	Ace:   12,
}

// SUIT represents one of the four suits
type SUIT string

const (
	// Clubs Suit constant
	Clubs = SUIT("C")
	// Diamonds Suit constant
	Diamonds = SUIT("D")
	// Hearts Suit constant
	Hearts = SUIT("H")
	// Spades Suit constant
	Spades = SUIT("S")
)

// Suits is a slice of all suits
var Suits = []SUIT{
	Clubs,
	Diamonds,
	Hearts,
	Spades,
}

// SuitIndexes ...
var SuitIndexes = map[SUIT]int{
	Clubs:    0,
	Diamonds: 1,
	Hearts:   2,
	Spades:   3,
}

// Card is the type that represents a playing card
type Card struct {
	Rank RANK
	Suit SUIT
}

// New creates a new instance of the card object
func New(v RANK, s SUIT) *Card {
	return &Card{
		Rank: v,
		Suit: s,
	}
}

// String representation of a card
func (c *Card) String() string {
	return fmt.Sprintf("%s%s", c.Rank, c.Suit)
}

// Connected returns true if the given
// card is adjacent (connected in a straight)
// to the current card, includes K <-> A, A <-> 2
func (c *Card) Connected(a *Card) bool {
	ours := RankIndexes[c.Rank]
	theirs := RankIndexes[a.Rank]
	// TODO: Handle Ace case more elegantly
	if ours == theirs+1 || ours == theirs-1 ||
		c.Rank == Ace && a.Rank == Two ||
		a.Rank == Ace && c.Rank == Two {
		return true
	}
	return false
}

func makeCardMatrix() [NumSuit][NumRanks]*Card {
	var matrix [NumSuit][NumRanks]*Card
	// populate cardMatrix
	for suitI := 0; suitI < NumSuit; suitI++ {
		for rankI := 0; rankI < NumRanks; rankI++ {
			matrix[suitI][rankI] = New(Ranks[rankI], Suits[suitI])
		}
	}
	return matrix
}
