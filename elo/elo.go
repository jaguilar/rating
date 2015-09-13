// Package elo implements the elo rating update rule.
package elo

import (
	"fmt"
	"math"

	"github.com/jaguilar/rating"
)

// System represents the elo rating system, including constants used for performing
// updates. It implements rating.System.
type System struct {
	// K is the K factor to be used for the update.
	K float64

	// RatingFloor is the minimum allowed rating.
	RatingFloor float64
}

func (s System) String() string {
	return fmt.Sprintf("%#v", s)
}

// WinChance is part of the rating.System interface. In elo, the win chance
// includes the half the chance of drawing, which is not specified.
func (s System) WinChance(r, ro rating.Rating) float64 {
	return expectedScore(r[0], ro[0])
}

// InitialRating is part of the rating.System interface.
func (s System) InitialRating() rating.Rating {
	return rating.Rating{800}
}

// Update is part of the rating.System interface.
func (s System) Update(r, ro rating.Rating, o rating.Outcome) rating.Rating {
	return rating.Rating{s.update(r[0], ro[0], score(o))}
}

func (s System) update(r, ro, score float64) float64 {
	e := expectedScore(r, ro)
	rn := r + s.K*(score-e)
	if rn < s.RatingFloor {
		rn = s.RatingFloor
	}
	return rn
}

func score(o rating.Outcome) float64 {
	wld := o.WLD()
	return wld.Value
}

// expectedScore returns the expected score a player with rating r would
// receive when playing against a player with a rating ro.
func expectedScore(r, ro float64) float64 {
	return 1. / (1. + math.Pow(10.0, (float64(ro)-float64(r))/400.0))
}
