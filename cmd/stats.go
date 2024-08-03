package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/dredly/ego/internal/db"
)

func RunStats() {
	statsCmd := flag.NewFlagSet("stats", flag.ExitOnError)
	name := statsCmd.String("name", "", "name of the player to display stats for")
	addVerboseFlag(statsCmd)
	addDbPathFlag(statsCmd)
	statsCmd.Parse(os.Args[2:])

	conn, err := db.Connect(dbPath, verbose)
	if err != nil {
		logger.Fatalf("failed to get db connection: %v", err)
	}

	gamesPlayed, err := conn.GamesPlayed(*name)
	if err != nil {
		logger.Fatalf("failed to determine number of games played by player: %v", err)
	}

	gamesWon, err := conn.GamesWon(*name)
	if err != nil {
		logger.Fatalf("failed to determine number of games won by player: %v", err)
	}

	peakELO, err := conn.PeakELOForPlayer(*name)
	if err != nil {
		logger.Fatalf("failed to determine peak ELO for player: %v", err)
	}

	displayedWinRate := "N/A"

	if gamesPlayed > 0 {
		winPercentage := (float64(gamesWon)  / float64(gamesPlayed)) * 100
		displayedWinRate = fmt.Sprintf("%.1f%%", winPercentage)
	}

	logger.Printf("Stats for %s", *name)
	logger.Printf("Games Played: %d", gamesPlayed)
	logger.Printf("Win rate: %s", displayedWinRate)
	logger.Printf("Peak ELO: %.4f", peakELO)
}