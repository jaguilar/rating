package elo

import (
	"testing"

	"github.com/jaguilar/rating"
	"github.com/jaguilar/testify/assert"
)

var _ rating.System = System{K: 32}

func TestElo(t *testing.T) {
	sys := System{K: 32}

	type opponentOutcome struct {
		ro float64
		o  rating.Outcome
	}

	for _, tc := range []struct {
		r, rf float64 // Initial, final rating.
		games []opponentOutcome
	}{
		{r: 1613, rf: 1603, games: []opponentOutcome{
			{ro: 1609, o: rating.Loss},
			{ro: 1477, o: rating.Draw},
			{ro: 1388, o: rating.Win},
			{ro: 1586, o: rating.Win},
			{ro: 1720, o: rating.Loss},
		}},
	} {
		r := tc.r
		for _, g := range tc.games {
			r = sys.Update(rating.Rating{r}, rating.Rating{g.ro}, g.o)[0]
		}
		assert.InDelta(t, float64(tc.rf), float64(r), .5)
	}
}
