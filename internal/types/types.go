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

func (p *Player) RecordWin(opponentELO float64, isDonut bool) {
	change := elo.EloChange(p.ELO, opponentELO, 1, isDonut)
	p.ELO += change
}

func (p *Player) RecordLoss(opponentELO float64, isDonut bool) {
	change := elo.EloChange(p.ELO, opponentELO, 0, isDonut)
	p.ELO += change
}

func (p *Player) RecordDraw(opponentELO float64) {
	change := elo.EloChange(p.ELO, opponentELO, 0.5, false)
	p.ELO += change
}