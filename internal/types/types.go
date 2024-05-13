package types

import "github.com/dredly/ego/internal/elo"

const startingELO = 1000

type Player struct {
	Name string
	ELO float64
}

func NewPlayer(name string) *Player {
	return &Player{
		Name: name,
		ELO: startingELO, 
	}
}

func (p *Player) RecordResult(opponentELO, score float64, multiplier int) {
	change := elo.EloChange(p.ELO, opponentELO, score)
	p.ELO += change * float64(multiplier)
}