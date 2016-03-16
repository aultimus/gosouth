package main

import (
	"fmt"

	"github.com/aultimus/gosouth/card"
	"github.com/aultimus/gosouth/hand"
	"github.com/aultimus/gosouth/headsup"
)

func main() {
	h1 := hand.Hand{
		card.New(card.Ace, card.Spades),
		card.New(card.Two, card.Spades)}
	h2 := hand.Hand{
		card.New(card.Nine, card.Clubs),
		card.New(card.Ten, card.Spades)}
	r, err := headsup.Prob(h1, h2)
	if err != nil {
		panic(err)
	}
	fmt.Println(r)
}
