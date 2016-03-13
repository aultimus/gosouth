package hand

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/aultimus/gosouth/card"
	"github.com/stretchr/testify/assert"
)

// TODO: test against a db of poker hands?
// csv?

// TODO: Test Sorting

func TestNumSuited(t *testing.T) {
	a := assert.New(t)
	h := Hand{
		card.New(card.Ace, card.Spades),
		card.New(card.Two, card.Spades),
		card.New(card.Nine, card.Spades),
		card.New(card.Four, card.Spades),
		card.New(card.Ace, card.Diamonds),
		card.New(card.Two, card.Hearts),
		card.New(card.Nine, card.Diamonds),
	}
	s, i := numSuited(h)
	a.Equal(card.Spades, s)
	a.Equal(4, i)
}

func TestRmDupsOfOtherSuits(t *testing.T) {
	a := assert.New(t)
	h := Hand{
		card.New(card.Ace, card.Spades),
		card.New(card.Two, card.Spades),
		card.New(card.Nine, card.Spades),
		card.New(card.Four, card.Spades),
		card.New(card.Ace, card.Diamonds),
		card.New(card.Two, card.Hearts),
		card.New(card.Nine, card.Diamonds),
	}
	e := Hand{
		card.New(card.Ace, card.Spades),
		card.New(card.Two, card.Spades),
		card.New(card.Nine, card.Spades),
		card.New(card.Four, card.Spades),
	}
	actual := rmDupsOfOtherSuits(h, card.Spades)
	a.True(reflect.DeepEqual(actual, e))
}

func TestStraight(t *testing.T) {
	a := assert.New(t)

	// Test wheel straight flush
	h := Hand{
		card.New(card.Ace, card.Clubs),
		card.New(card.Two, card.Clubs),
		card.New(card.Three, card.Clubs),
		card.New(card.Four, card.Clubs),
		card.New(card.Five, card.Clubs),
		card.New(card.Two, card.Diamonds),
		card.New(card.Nine, card.Diamonds),
	}
	b, f, formedHand := straight(h)
	a.True(b)
	a.True(f)
	a.Equal(card.New(card.Five, card.Clubs), formedHand[sizeHand-1])

	// Test for straight flush when theres a flush
	// and a straight but no straight flush
	h = Hand{
		card.New(card.Ace, card.Clubs),
		card.New(card.Two, card.Clubs),
		card.New(card.Three, card.Hearts),
		card.New(card.Four, card.Hearts),
		card.New(card.Five, card.Clubs),
		card.New(card.Jack, card.Diamonds),
		card.New(card.Nine, card.Clubs),
	}
	b, f, formedHand = straight(h)
	a.True(b)
	a.False(f)
	a.Equal(card.New(card.Five, card.Clubs), formedHand[sizeHand-1])

	h = Hand{
		card.New(card.Ten, card.Hearts),
		card.New(card.Nine, card.Hearts),
		card.New(card.King, card.Spades),
		card.New(card.Jack, card.Hearts),
		card.New(card.Queen, card.Diamonds),
		card.New(card.Eight, card.Hearts),
		card.New(card.Nine, card.Spades),
	}
	b, f, formedHand = straight(h)
	a.True(b)
	a.False(f)
	a.Equal(card.New(card.King, card.Spades), formedHand[sizeHand-1])

	h = Hand{
		card.New(card.Ten, card.Clubs),
		card.New(card.Two, card.Clubs),
		card.New(card.King, card.Clubs),
		card.New(card.Jack, card.Clubs),
		card.New(card.Queen, card.Clubs),
		card.New(card.Queen, card.Diamonds),
		card.New(card.Eight, card.Diamonds),
	}
	b, f, formedHand = straight(h)
	a.False(b)

	// TODO: Check a hand of < 5 cards
}

func TestStraightTieBreak(t *testing.T) {
	a := assert.New(t)
	holeHand1 := Hand{
		card.New(card.Jack, card.Hearts),
		card.New(card.King, card.Hearts),
	}

	holeHand2 := Hand{
		card.New(card.Six, card.Hearts),
		card.New(card.Nine, card.Spades),
	}
	commCards := Hand{
		card.New(card.Ten, card.Hearts),
		card.New(card.Nine, card.Hearts),
		card.New(card.Eight, card.Hearts),
		card.New(card.Seven, card.Hearts),
		card.New(card.Ace, card.Spades),
	}
	h1 := append(holeHand1, commCards...)
	h2 := append(holeHand2, commCards...)

	hasStraight, hasFlush, formedHand1 := straight(h1)
	a.True(hasStraight, fmt.Sprintf("%s does not have a straight?", h1))
	a.True(hasFlush)
	a.Equal(card.New(card.Jack, card.Hearts), formedHand1[sizeHand-1])

	hasStraight, hasFlush, formedHand2 := straight(h2)
	a.True(hasStraight, fmt.Sprintf("%s does not have a straight?", h2))
	a.True(hasFlush)
	a.Equal(card.New(card.Ten, card.Hearts), formedHand2[sizeHand-1])

	v1 := NewHandValue(StraightFlush, formedHand1)
	v2 := NewHandValue(StraightFlush, formedHand2)

	outcome := tieBreak(v1, v2)
	a.Equal(H1Win, outcome)
}

func TestFlush(t *testing.T) {
	a := assert.New(t)
	// basic test
	h := Hand{
		card.New(card.Four, card.Clubs),
		card.New(card.Two, card.Clubs),
		card.New(card.King, card.Clubs),
		card.New(card.Jack, card.Hearts),
		card.New(card.Queen, card.Clubs),
		card.New(card.Queen, card.Spades),
		card.New(card.Eight, card.Clubs),
	}
	f, h := flush(h)
	a.True(f)
	a.True(reflect.DeepEqual(Hand{
		card.New(card.Two, card.Clubs),
		card.New(card.Four, card.Clubs),
		card.New(card.Eight, card.Clubs),
		card.New(card.Queen, card.Clubs),
		card.New(card.King, card.Clubs),
	}, h))

	// test with 6 of the same suit
	h = Hand{
		card.New(card.Four, card.Clubs),
		card.New(card.Three, card.Clubs),
		card.New(card.King, card.Clubs),
		card.New(card.Ace, card.Clubs),
		card.New(card.Queen, card.Clubs),
		card.New(card.Queen, card.Spades),
		card.New(card.Eight, card.Clubs),
	}
	f, h = flush(h)
	a.True(f)
	a.True(reflect.DeepEqual(Hand{
		card.New(card.Four, card.Clubs),
		card.New(card.Eight, card.Clubs),
		card.New(card.Queen, card.Clubs),
		card.New(card.King, card.Clubs),
		card.New(card.Ace, card.Clubs),
	}, h))

	h = Hand{
		card.New(card.Four, card.Spades),
		card.New(card.Ace, card.Spades),
		card.New(card.King, card.Hearts),
		card.New(card.Jack, card.Hearts),
		card.New(card.Queen, card.Clubs),
		card.New(card.Queen, card.Spades),
		card.New(card.Eight, card.Clubs),
	}
	f, h = flush(h)
	a.False(f)
}

func TestFourOfAKind(t *testing.T) {
	a := assert.New(t)
	h := Hand{
		card.New(card.Jack, card.Diamonds),
		card.New(card.Jack, card.Spades),
		card.New(card.King, card.Clubs),
		card.New(card.Jack, card.Hearts),
		card.New(card.Four, card.Clubs),
		card.New(card.Three, card.Spades),
		card.New(card.Jack, card.Clubs),
	}
	b, h := xOfAKind(h, 4)
	a.True(b)
	a.Equal(card.New(card.King, card.Clubs), h[sizeHand-1])
	a.Equal(card.Jack, h[0].Rank)

	h = Hand{
		card.New(card.Four, card.Spades),
		card.New(card.Ace, card.Spades),
		card.New(card.Queen, card.Hearts),
		card.New(card.Jack, card.Hearts),
		card.New(card.Queen, card.Clubs),
		card.New(card.Queen, card.Spades),
		card.New(card.Eight, card.Clubs),
	}
	b, _ = xOfAKind(h, 4)
	a.False(b)
}

func TestThreeOfAKind(t *testing.T) {
	a := assert.New(t)
	h := Hand{
		card.New(card.Jack, card.Diamonds),
		card.New(card.Jack, card.Spades),
		card.New(card.King, card.Clubs),
		card.New(card.Jack, card.Hearts),
		card.New(card.Four, card.Clubs),
		card.New(card.Three, card.Spades),
		card.New(card.Ace, card.Clubs),
	}
	b, h := xOfAKind(h, 3)
	a.True(b)
	a.Equal(card.Jack, h[0].Rank)                           // rank of the set
	a.Equal(card.New(card.Ace, card.Clubs), h[sizeHand-2])  // first kicker
	a.Equal(card.New(card.King, card.Clubs), h[sizeHand-1]) // second kicker

	h = Hand{
		card.New(card.Four, card.Spades),
		card.New(card.Ace, card.Spades),
		card.New(card.King, card.Hearts),
		card.New(card.Jack, card.Hearts),
		card.New(card.Queen, card.Clubs),
		card.New(card.Queen, card.Spades),
		card.New(card.Eight, card.Clubs),
	}
	b, _ = xOfAKind(h, 3)
	a.False(b)

	h = Hand{
		card.New(card.Ace, card.Spades),
		card.New(card.King, card.Spades),
		card.New(card.Ace, card.Hearts),
		card.New(card.King, card.Hearts),
		card.New(card.King, card.Clubs),
		card.New(card.Ace, card.Clubs),
		card.New(card.Two, card.Diamonds),
	}
	b, h = xOfAKind(h, 3)
	a.True(b)
	a.Equal(card.Ace, h[0].Rank)           // rank of the set
	a.Equal(card.King, h[sizeHand-2].Rank) // first kicker
	a.Equal(card.King, h[sizeHand-1].Rank) // second kicker

	h = Hand{
		card.New(card.Queen, card.Spades),
		card.New(card.King, card.Spades),
		card.New(card.Queen, card.Hearts),
		card.New(card.King, card.Hearts),
		card.New(card.King, card.Clubs),
		card.New(card.Queen, card.Clubs),
		card.New(card.Ace, card.Diamonds),
	}
	b, h = xOfAKind(h, 3)
	a.True(b)
	a.Equal(card.King, h[0].Rank)           // rank of the set
	a.Equal(card.Ace, h[sizeHand-2].Rank)   // first kicker
	a.Equal(card.Queen, h[sizeHand-1].Rank) // second kicker
}

func TestPopKicker(t *testing.T) {
	a := assert.New(t)

	h := Hand{
		card.New(card.Queen, card.Spades),
		card.New(card.King, card.Spades),
		card.New(card.Queen, card.Hearts),
		card.New(card.King, card.Hearts),
		card.New(card.King, card.Clubs),
		card.New(card.Queen, card.Clubs),
		card.New(card.Ace, card.Diamonds),
	}
	f, k := popKicker(h)
	a.Equal(card.New(card.Ace, card.Diamonds), k)
	a.True(reflect.DeepEqual(f, Hand{card.New(card.Queen, card.Spades),
		card.New(card.King, card.Spades),
		card.New(card.Queen, card.Hearts),
		card.New(card.King, card.Hearts),
		card.New(card.King, card.Clubs),
		card.New(card.Queen, card.Clubs),
	}))
}

func TestFullHouse(t *testing.T) {
	a := assert.New(t)
	h := Hand{
		card.New(card.Queen, card.Spades),
		card.New(card.King, card.Spades),
		card.New(card.Queen, card.Hearts),
		card.New(card.King, card.Hearts),
		card.New(card.Queen, card.Clubs),
		card.New(card.Two, card.Clubs),
		card.New(card.Ace, card.Diamonds),
	}
	b, f := fullHouse(h)
	a.True(b)
	a.Equal(card.Queen, f[0].Rank)
	a.Equal(card.Queen, f[1].Rank)
	a.Equal(card.Queen, f[2].Rank)
	a.Equal(card.King, f[3].Rank)
	a.Equal(card.King, f[4].Rank)

	h = Hand{
		card.New(card.Queen, card.Spades),
		card.New(card.King, card.Spades),
		card.New(card.Ace, card.Hearts),
		card.New(card.King, card.Hearts),
		card.New(card.Queen, card.Clubs),
		card.New(card.Two, card.Clubs),
		card.New(card.Ace, card.Diamonds),
	}
	b, f = fullHouse(h)
	a.False(b)

	h = Hand{
		card.New(card.Two, card.Spades),
		card.New(card.Ace, card.Spades),
		card.New(card.Ace, card.Hearts),
		card.New(card.King, card.Hearts),
		card.New(card.Queen, card.Clubs),
		card.New(card.Seven, card.Clubs),
		card.New(card.Ace, card.Diamonds),
	}
	b, f = fullHouse(h)
	a.False(b)
}

func TestTwoPair(t *testing.T) {
	a := assert.New(t)
	h := Hand{
		card.New(card.Queen, card.Spades),
		card.New(card.King, card.Spades),
		card.New(card.Queen, card.Hearts),
		card.New(card.King, card.Hearts),
		card.New(card.Seven, card.Clubs),
		card.New(card.Two, card.Clubs),
		card.New(card.Ace, card.Diamonds),
	}
	b, f := twoPair(h)
	a.True(b)
	a.Equal(card.King, f[0].Rank)
	a.Equal(card.King, f[1].Rank)
	a.Equal(card.Queen, f[2].Rank)
	a.Equal(card.Queen, f[3].Rank)
	a.Equal(card.Ace, f[4].Rank)

	h = Hand{
		card.New(card.Nine, card.Spades),
		card.New(card.King, card.Spades),
		card.New(card.Ace, card.Hearts),
		card.New(card.Two, card.Hearts),
		card.New(card.Queen, card.Clubs),
		card.New(card.Two, card.Clubs),
		card.New(card.Seven, card.Diamonds),
	}
	b, f = twoPair(h)
	a.False(b)
}

func TestOnePair(t *testing.T) {
	a := assert.New(t)
	h := Hand{
		card.New(card.Queen, card.Spades),
		card.New(card.King, card.Spades),
		card.New(card.Queen, card.Hearts),
		card.New(card.Three, card.Hearts),
		card.New(card.Seven, card.Clubs),
		card.New(card.Two, card.Clubs),
		card.New(card.Ace, card.Diamonds),
	}
	b, f := onePair(h)
	a.True(b)
	a.Equal(card.Queen, f[0].Rank)
	a.Equal(card.Queen, f[1].Rank)
	a.Equal(card.Ace, f[2].Rank)
	a.Equal(card.King, f[3].Rank)
	a.Equal(card.Seven, f[4].Rank)

	h = Hand{
		card.New(card.Nine, card.Spades),
		card.New(card.King, card.Spades),
		card.New(card.Ace, card.Hearts),
		card.New(card.Two, card.Hearts),
		card.New(card.Queen, card.Clubs),
		card.New(card.Five, card.Clubs),
		card.New(card.Seven, card.Diamonds),
	}
	b, f = onePair(h)
	a.False(b)
}
