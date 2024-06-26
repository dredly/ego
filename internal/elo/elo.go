package elo

import "math"

const kFactor float64 = 20

func EloChange(r1, r2, score float64) float64 {
	return kFactor * (score - expectedScore(r1, r2))
}

func expectedScore(r1, r2 float64) float64 {
	return 1 / (1 + math.Pow(float64(10), (r2 - r1)/400))
}