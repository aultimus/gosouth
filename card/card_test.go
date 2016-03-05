package card

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var cardMatrix = makeCardMatrix()

func TestMain(m *testing.M) {

	ret := m.Run()
	os.Exit(ret)
}

// Test the connectedness of every card with every other card
func TestConnected(t *testing.T) {
	a := assert.New(t)
	// suit loop
	for suitI := 0; suitI < NumSuit; suitI++ {
		// rank loop
		for rankI := 0; rankI < NumRanks; rankI++ {
			c := cardMatrix[suitI][rankI]
			// other suit loop
			for otherSuitI := 0; otherSuitI < NumSuit; otherSuitI++ {
				// other rank loop
				for otherRankI := 0; otherRankI < NumRanks; otherRankI++ {
					o := cardMatrix[otherSuitI][otherRankI]
					rankDiff := RankIndexes[c.Rank] - RankIndexes[o.Rank]
					// TODO: Nicer comp for Ace <-> Two
					if rankDiff == 1 || rankDiff == -1 ||
						c.Rank == Ace && o.Rank == Two ||
						o.Rank == Ace && c.Rank == Two {
						a.True(c.Connected(o), fmt.Sprintf("%s should be connected with %s", c, o))
					} else {
						a.False(c.Connected(o), fmt.Sprintf("%s should not be connected with %s", c, o))
					}
				}
			}
		}
	}
}
