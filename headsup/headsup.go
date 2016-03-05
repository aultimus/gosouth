package headsup

import (
	"fmt"

	"github.com/aultimus/gosouth/deck"
)

// Result represents the probability breakdown
// of a hand unfolding
type Result struct {
	H1Win float64
	H2Win float64
	Tie   float64
}

// Prob given two initial starting hands
// calculates the probabilities of the results
// by simulating every possible deal from the
// resultant deck
func Prob(h1, h2 deck.Deck) (*Result, error) {
	var r *Result
	var err error
	c := make(chan deck.Deck)
	d := deck.New()
	for _, v := range append(h1, h2...) {
		d, err = deck.Remove(d, v)
		if err != nil {
			return r, err
		}
	}
	go deck.Combs(d, 5, c)
	for {
		select {
		case v := <-c:
			// Todo: generate all possible deals
			// and see whats the best hand that can be
			// made using FormHand() and count showdown
			// results out of Showdown()
			fmt.Println(v)
		}
	}
}
