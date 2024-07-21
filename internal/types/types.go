package types

import (
	"time"

	"github.com/dredly/ego/internal/elo"
)


type Player struct {
	ID int
	Name string
	ELO float64
}

func NewPlayer(name string, ELO float64) *Player {
	return &Player{
		Name: name,
		ELO: ELO, 
	}
}

func (p *Player) RecordResult(opponentELO, score float64, multiplier int) {
	change := elo.EloChange(p.ELO, opponentELO, score)
	p.ELO += change * float64(multiplier)
}

type Game struct {
	ID, Player1ID, Player2ID, Player1Points, Player2Points int
	Player1ELOBefore, Player2ELOBefore, Player1ELOAfter, Player2ELOAfter float64
	Played time.Time
}

type GameDisplay struct {
	Player1Name, Player2Name string
	Player1Points, Player2Points int
	Played time.Time
}