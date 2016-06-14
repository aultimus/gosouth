package headsup

import (
	"fmt"

	"github.com/aultimus/gosouth/deck"
	"github.com/aultimus/gosouth/hand"
)

// Result represents the probability breakdown
// of a hand unfolding
type Result struct {
	Win []float64
}

// NewResult creates a new Result instance
// for comparing numHands hands
func NewResult(numHands int) *Result {
	r := Result{}
	r.Win = make([]float64, numHands, numHands)
	return &r
}

func (r *Result) String() string {
	var s string
	for i, w := range r.Win {
		s = fmt.Sprintf("%s H%d: %0.2f", s, i, w)
	}
	return s
}

// Prob given n initial starting hands calculates the probabilities of the
// results by simulating every possible deal from the resultant deck.
// Likely faster to use a lookup table, this function can help generate one
func Prob(hands ...hand.Hand) (*Result, error) {
	r := NewResult(len(hands))
	var err error
	c := make(chan deck.Deck)
	d := deck.New()

	var usedCards hand.Hand
	for _, h := range hands {
		usedCards = append(usedCards, h...)
	}

	d, err = deck.RemoveMultiple(d, usedCards)
	if err != nil {
		return r, err
	}
	go deck.Combs(d, 5, c)
	count := 0
	for v := range c {
		count++
		var pHands []hand.Hand
		for _, h := range hands {
			pHands = append(pHands, hand.Hand(append(v, h...)))
		}

		winners := hand.Showdown(pHands)
		// draws will add up to over 100% but we are ok with that
		for i := 0; i < len(winners); i++ {
			r.Win[i]++
		}
	}

	// convert to percentages
	var total float64
	for _, v := range r.Win {
		total += v
	}
	for i := 0; i < len(r.Win); i++ {
		r.Win[i] = r.Win[i] / total * 100
	}
	return r, nil
}

// TODO: Write a function that will generate a
// table of all possible matchups
