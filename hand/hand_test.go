package hand

import (
	"reflect"
	"testing"

	"github.com/aultimus/gosouth/card"
	"github.com/stretchr/testify/assert"
)

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
	b, r,f = straight(h)
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
	b, r,f  = straight(h)
	a.False(b)
	a.Equal(card.Nil, r)
}
