package deck

import (
	"reflect"
	"testing"

	"github.com/aultimus/gosouth/card"
	"github.com/stretchr/testify/assert"
)

func TestShuffle(t *testing.T) {
	a := assert.New(t)
	SeedWithNow()
	d := New()
	for i := 0; i < 1000; i++ {
		c := NewShuffled()
		a.Equal(false, reflect.DeepEqual(d, c))
	}
}

func TestRemove(t *testing.T) {
	a := assert.New(t)
	var err error
	d := New()
	d, err = Remove(d, card.New(card.Ace, card.Spades))
	a.NoError(err)
	d, err = Remove(d, card.New(card.Ace, card.Spades))
	a.Error(err)
}
