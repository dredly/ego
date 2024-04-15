package elo

import (
	"math"
	"testing"
)

func Test_EloChange(t *testing.T) {
	eloChange := EloChange(1750, 1810, 1)
	expectedEloChange := 11.7
	if !equalWithinTolerance(t, eloChange, expectedEloChange, 0.1) {
		t.Errorf("got eloChange of %f but wanted %f", eloChange, expectedEloChange)
	}
}

func Test_expectedScore_equal_rankings(t *testing.T) {
	expectedScore := expectedScore(1000, 1000)
	if expectedScore != 0.5 {
		t.Errorf("expectedScore should have been 0.5 but was %f", expectedScore)
	}
}

// TODO: table based test in this format
func Test_expectedScore(t *testing.T) {
	gotExpectedScore := expectedScore(1750, 1810)
	wantExpectedScore := 0.414
	if !equalWithinTolerance(t, gotExpectedScore, wantExpectedScore, 0.01) {
		t.Errorf("got expectedScore of %f but wanted %f", gotExpectedScore, wantExpectedScore)
	}
}

func equalWithinTolerance(t *testing.T, a, b, tolerance float64) bool {
	t.Helper()
	diff := math.Abs(a - b)
	return diff <= tolerance
}