package elo

import (
	"testing"

	"github.com/jaguilar/testify/assert"
)

func TestElo(t *testing.T) {
	config := Config{K: 32.0}

	type opponentOutcome struct {
		ro Rating
		o  Outcome
	}

	for _, tc := range []struct {
		r, rf Rating // Initial, final rating.
		games []opponentOutcome
	}{
		{r: 1613, rf: 1601, games: []opponentOutcome{
			{ro: 1609, o: Loss},
			{ro: 1477, o: Draw},
			{ro: 1388, o: Win},
			{ro: 1586, o: Win},
			{ro: 1720, o: Loss},
		}},
	} {
		r := tc.r
		for _, g := range tc.games {
			r = Update(r, g.ro, g.o, config)
		}
		assert.Equal(t, tc.rf, r)
	}
}
