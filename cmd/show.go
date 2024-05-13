package cmd

import (
	"flag"

	"github.com/dredly/ego/internal/db"
)

func RunShow() {
	showCmd := flag.NewFlagSet("show", flag.ExitOnError)
	addVerboseFlag(showCmd)
	conn, err := db.New()
	if err != nil {
		logger.Fatalf("failed to get db connection: %v", err)
	}
	players, err:= conn.Show()
	if err != nil {
		logger.Fatalf("failed to show all players: %v", err)
	}
	logger.Printf("Leaderboard")
	for i, player := range players {
		logger.Printf("%d. %s: %.2f", i + 1, player.Name, player.ELO)
	}
}