package deck

import (
	"fmt"
	"math/rand"
	"reflect"
	"time"

	"github.com/aultimus/gosouth/card"
)

// Deck represents a deck of cards
type Deck []*card.Card

// New returns a fresh unsorted fifty-two card deck
func New() Deck {
	var d []*card.Card
	for _, s := range card.Suits {
		for _, v := range card.Ranks {
			d = append(d, card.New(v, s))
		}
	}
	return d
}

// NewShuffled returns a sorted deck
func NewShuffled() Deck {
	return knuthShuffle(New())
}

// SeedWithNow seeds the random number gen
// with the current time (millis)
func SeedWithNow() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// Combs returns all possible combinations of
// dealing k cards from Deck d
// from len(d) pick k
func Combs(d Deck, k int, c chan Deck) {
	pool := d
	n := len(pool)

	if k > n {
		return
	}

	indices := make([]int, k)
	for i := range indices {
		indices[i] = i
	}

	result := make(Deck, k)
	for i, el := range indices {
		result[i] = pool[el]
	}
	c <- result

	for {
		i := k - 1
		for ; i >= 0 && indices[i] == i+n-k; i-- {
		}
		if i < 0 {
			close(c)
			return
		}
		indices[i]++
		for j := i + 1; j < k; j++ {
			indices[j] = indices[j-1] + 1
		}

		for ; i < len(indices); i++ {
			result[i] = pool[indices[i]]
		}
		c <- result
	}
}

// Remove removes a card from the given deck,
// returns an error if given card is not in deck
func Remove(d Deck, c *card.Card) (Deck, error) {
	for i, v := range d {
		if reflect.DeepEqual(c, v) {
			// Garbage collection probs?
			// https://github.com/golang/go/wiki/SliceTricks
			d = append(d[:i], d[i+1:]...)
			return d, nil
		}
	}
	return d, fmt.Errorf("card %s is not in deck", c)
}

// knuthShuffle is an implementation of the
// Knuth/Fisher-Yates shuffle, in place, O(n)
// https://en.wikipedia.org/wiki/Fisher%E2%80%93Yates_shuffle
// Remember to seed prior to calling
func knuthShuffle(a Deck) Deck {
	n := len(a)
	for i := 0; i < n-2; i++ {
		j := rand.Intn(n - i)
		a[i], a[i+j] = a[i+j], a[i]
	}
	return a
}
