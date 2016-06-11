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
	Tie float64
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
	s = fmt.Sprintf("%s Draw: %0.2f", s, r.Tie)
	return s
}

// TODO: Add Prob(h hand.Hand) (*Result, error) func
// that calculates pre-flop prob of a given hand being the best hand

// Prob given two initial starting hands
// calculates the probabilities of the results
// by simulating every possible deal from the
// resultant deck
func Prob(h1, h2 hand.Hand) (*Result, error) {
	r := NewResult(2)
	var err error
	c := make(chan deck.Deck)
	d := deck.New()
	d, err = deck.RemoveMultiple(d, append(h1, h2...))
	if err != nil {
		return r, err
	}
	go deck.Combs(d, 5, c)
	count := 0
	for v := range c {
		count++
		p1Hand := hand.Hand(append(v, h1...))
		p2Hand := hand.Hand(append(v, h2...))
		outcome := hand.Showdown(p1Hand, p2Hand)
		if outcome == hand.H1Win {
			r.Win[0]++
		} else if outcome == hand.H2Win {
			r.Win[1]++
		} else {
			r.Tie++
		}
	}
	total := r.Win[0] + r.Win[1] + r.Tie
	r.Win[0] = r.Win[0] / total * 100
	r.Win[1] = r.Win[1] / total * 100
	r.Tie = r.Tie / total * 100
	return r, nil
}

// TODO: Write a function that will generate a
// table of all possible matchups
