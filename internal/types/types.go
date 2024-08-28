package types

import (
	"fmt"
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

type GameRecording struct {
	Player1 PlayerRecording
	Player2 PlayerRecording
}

type PlayerRecording struct {
	Player Player
	Points int
	ELOBefore float64
	ELOAfter float64
}

type Game struct {
	ID, Player1ID, Player2ID, Player1Points, Player2Points int
	Player1ELOBefore, Player2ELOBefore, Player1ELOAfter, Player2ELOAfter float64
	Played time.Time
}

type GameSummary struct {
	Player1Name, Player2Name string
	Player1Points, Player2Points int
	Played time.Time
}

type GameDetail struct {
	GameSummary
	Player1ELOBefore, Player2ELOBefore, Player1ELOAfter, Player2ELOAfter float64
	ID int
}

type GameResults struct {
	Won, Drawn, Lost int
}

func (gr GameResults) Total() int {
	return gr.Won + gr.Drawn + gr.Lost
}

func (gr GameResults) WinRateStr() string {
	if gr.Total() == 0 {
		return "N/A"
	}
	winPercentage := (float64(gr.Won)  / float64(gr.Total())) * 100
	return fmt.Sprintf("%.1f%%", winPercentage)
}