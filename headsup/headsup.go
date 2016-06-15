package headsup

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/aultimus/gosouth/deck"
	"github.com/aultimus/gosouth/hand"
)

// TODO: Rename this package 'prob'

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

// Prob given n (1 -> many) initial starting hands calculates the probabilities
// of the results by simulating every possible deal from the resultant deck.
// Likely faster to use a lookup table, this function can help generate one
func Prob(hands ...hand.Hand) (*Result, error) {
	numResults := len(hands)
	if len(hands) == 1 {
		numResults = 2
	}
	r := NewResult(numResults)
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
	numCardsToDeal := 5
	if len(hands) == 1 {
		numCardsToDeal = 7
	}
	go deck.Combs(d, numCardsToDeal, c)
	count := 0
	var extraHands []hand.Hand
	for v := range c {
		count++
		var pHands []hand.Hand
		// generate an opposing hand out of the deal
		if len(hands) == 1 {
			extraHands = []hand.Hand{{v[0], v[1]}}
			v = v[2:]
		}

		for _, h := range append(hands, extraHands...) {
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

func rmWhitespace(s string) string {
	return strings.Replace(s, " ", "", -1)
}

// build a map of hand type (there being 169 of such) rather than 2,652 possible
// starting hands to percentage of winning.

// HandProb represents the probability of a hand winning
type HandProb struct {
	Win float64
	Tie float64
}

// HandProbMap returns a map of hand types (string of format of one of
// {89o, 89s, 99}.
// TODO: Add test coverage
func HandProbMap() map[string]HandProb {
	f, err := os.Open("hand_types.csv")
	if err != nil {
		panic(err.Error())
	}
	r := csv.NewReader(f)
	r.Comment = '#'
	lines, err := r.ReadAll()
	if err != nil {
		panic(err.Error())
	}
	f.Close()

	// note: cannot use structs as map keys
	handSuccessMap := make(map[string]HandProb)
	for _, line := range lines {
		handType := rmWhitespace(line[0])
		v := rmWhitespace(line[1])
		fmt.Println(handType)
		fmt.Println(v)
		winPer, err := strconv.ParseFloat(v, 32)
		if err != nil {
			panic(err)
		}
		tiePer, err := strconv.ParseFloat(rmWhitespace(line[2]), 32)
		if err != nil {
			panic(err)
		}
		h := HandProb{Win: winPer, Tie: tiePer}
		handSuccessMap[handType] = h
	}
	return handSuccessMap
}
