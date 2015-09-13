// Package elo implements the elo rating update rule.
package elo

import "math"

// Rating is the rating of a player or team.
type Rating float64

// Outcome is the result of a match.
type Outcome float64

// Possible outcomes.
const (
	Win            Outcome = 1.0
	Loss           Outcome = 0.0
	Draw           Outcome = 0.5
	UnknownOutcome Outcome = -1
)

// Opposite returns the opposite outcome from the receiver.
func (o Outcome) Opposite() Outcome {
	switch o {
	case Win, Loss:
		return Outcome(1.0 - float64(o))
	default:
		return o
	}
}

// ParseOutcome returns the outcome represented by s, or else UnknownOutcome.
func ParseOutcome(s string) Outcome {
	switch s {
	case "win":
		return Win
	case "loss":
		return Loss
	case "draw":
		return Draw
	default:
		return UnknownOutcome
	}
}

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
