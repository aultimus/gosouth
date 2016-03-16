package deck

import (
	"fmt"
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

// Long running test
func TestCombs(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping TestCombs as in short mode")
	}
	a := assert.New(t)
	c := make(chan Deck)
	d := New()
	go Combs(d, 5, c)
	total := 2598960 // choose 5 from 52
	count := 0
	for {
		select {
		case v, ok := <-c:
			if !ok {
				a.Equal(total, count)
				return
			}
			count++
			a.Equal(5, len(v), fmt.Sprintf("%s does not have 5 cards", v))
		}
	}
}
