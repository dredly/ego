package cmd

import (
	"flag"
	"os"

	"github.com/dredly/ego/internal/db"
)

func RunLeaderboard() {
	leaderboardCmd := flag.NewFlagSet("leaderboard", flag.ExitOnError)
	addVerboseFlag(leaderboardCmd)
	addDbPathFlag(leaderboardCmd)
	leaderboardCmd.Parse(os.Args[2:])

	conn, err := db.Connect(dbPath, verbose)
	if err != nil {
		logger.Fatalf("failed to get db connection: %v", err)
	}

	players, err := conn.AllPlayers()
	if err != nil {
		logger.Fatalf("failed to show all players: %v", err)
	}
	logger.Printf("Leaderboard")
	for i, player := range players {
		logger.Printf("%d. %s: %.2f", i + 1, player.Name, player.ELO)
	}
}