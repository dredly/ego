package types

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

func (p *Player) UpdateELO(change float64) {	
	p.ELO += change
}