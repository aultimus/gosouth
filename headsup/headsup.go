package headsup

import (
	"github.com/aultimus/gosouth/deck"
	"github.com/aultimus/gosouth/hand"
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
func Prob(h1, h2 hand.Hand) (*Result, error) {
	r := &Result{}
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
	count := 0
	for {
		select {
		case v, ok := <-c:
			if !ok {
				goto DONE
			}
			count++
			p1Hand := hand.Hand(append(v, h1...))
			p2Hand := hand.Hand(append(v, h2...))
			outcome := hand.Showdown(p1Hand, p2Hand)
			if outcome == hand.H1Win {
				r.H1Win++
			} else if outcome == hand.H2Win {
				r.H2Win++
			} else {
				r.Tie++
			}
		}
	}
DONE:
	total := r.H1Win + r.H2Win + r.Tie
	r.H1Win = r.H1Win / total * 100
	r.H2Win = r.H2Win / total * 100
	r.Tie = r.Tie / total * 100
	return r, nil
}
