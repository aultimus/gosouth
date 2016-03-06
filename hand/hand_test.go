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

	h := Hand{
		card.New(card.Ace, card.Clubs),
		card.New(card.Two, card.Clubs),
		card.New(card.Three, card.Clubs),
		card.New(card.Four, card.Clubs),
		card.New(card.Five, card.Clubs),
		card.New(card.Two, card.Diamonds),
		card.New(card.Nine, card.Diamonds),
	}
	b, r, f := straight(h)
	a.True(b)
	a.True(f)
	a.Equal(card.Five, r)

	h = Hand{
		card.New(card.Ten, card.Hearts),
		card.New(card.Nine, card.Hearts),
		card.New(card.King, card.Spades),
		card.New(card.Jack, card.Hearts),
		card.New(card.Queen, card.Diamonds),
		card.New(card.Eight, card.Hearts),
		card.New(card.Nine, card.Spades),
	}
	b, r, f = straight(h)
	a.True(b)
	a.False(f)
	a.Equal(card.King, r)

	h = Hand{
		card.New(card.Ten, card.Clubs),
		card.New(card.Two, card.Clubs),
		card.New(card.King, card.Clubs),
		card.New(card.Jack, card.Clubs),
		card.New(card.Queen, card.Clubs),
		card.New(card.Queen, card.Diamonds),
		card.New(card.Eight, card.Diamonds),
	}
	b, r, f = straight(h)
	a.False(b)
	a.Equal(card.Nil, r)

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

	hasStraight, rankValue1, _ := straight(h1)
	a.True(hasStraight, fmt.Sprintf("%s does not have a straight?", h1))

	hasStraight, rankValue2, _ := straight(h2)
	a.True(hasStraight, fmt.Sprintf("%s does not have a straight?", h2))

	// TODO: This code is duplicated
	v1 := NewHandValue(Straight, card.Nil, rankValue1, h1)
	v2 := NewHandValue(Straight, card.Nil, rankValue2, h2)

	outcome, err := straightTieBreak(v1, v2)
	a.NoError(err)
	a.Equal(H2Win, outcome)
}

func TestFlush(t *testing.T) {
	a := assert.New(t)
	h := Hand{
		card.New(card.Four, card.Clubs),
		card.New(card.Two, card.Clubs),
		card.New(card.King, card.Clubs),
		card.New(card.Jack, card.Hearts),
		card.New(card.Queen, card.Clubs),
		card.New(card.Queen, card.Spades),
		card.New(card.Eight, card.Clubs),
	}
	f := flush(h)
	a.True(f)

	h = Hand{
		card.New(card.Four, card.Spades),
		card.New(card.Ace, card.Spades),
		card.New(card.King, card.Hearts),
		card.New(card.Jack, card.Hearts),
		card.New(card.Queen, card.Clubs),
		card.New(card.Queen, card.Spades),
		card.New(card.Eight, card.Clubs),
	}
	f = flush(h)
	a.False(f)
}
