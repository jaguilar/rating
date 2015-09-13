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
		{r: 1613, rf: 1603, games: []opponentOutcome{
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
		assert.InDelta(t, float64(tc.rf), float64(r), .5)
	}
}

func TestOpposite(t *testing.T) {
	assert.Equal(t, Loss, Win.Opposite())
	assert.Equal(t, Win, Loss.Opposite())
	assert.Equal(t, Draw, Draw.Opposite())
}
