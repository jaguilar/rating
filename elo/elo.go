// Package elo implements the elo rating update rule.
package elo

import "math"

// Rating is the rating of a player or team.
type Rating float64

// Outcome is the result of a match.
type Outcome float64

// Possible outcomes.
const (
	Win  Outcome = 1.0
	Loss Outcome = 0.0
	Draw Outcome = 0.5
)

// Config contains the parameters of the update.
type Config struct {
	// K is the K factor to be used for the update.
	K float64

	// RatingFloor is the minimum allowed rating.
	RatingFloor float64
}

// Update returns an updated rating for a player with rating r playing against
// a player with a rating ro.
func Update(r, ro Rating, o Outcome, c Config) Rating {
	e := expectedScore(r, ro)
	rn := float64(r) + c.K*(float64(o)-e)
	if rn < c.RatingFloor {
		rn = c.RatingFloor
	}
	return Rating(rn)
}

// expectedScore returns the expected score a player with rating r would
// receive when playing against a player with a rating ro.
func expectedScore(r, ro Rating) float64 {
	return 1. / (1. + math.Pow(10.0, (float64(ro)-float64(r))/400.0))
}
