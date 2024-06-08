package cmd

import (
	"flag"
	"os"

	"github.com/dredly/ego/internal/db"
)

func RunGames() {
	gamesCmd := flag.NewFlagSet("games", flag.ExitOnError)
	addVerboseFlag(gamesCmd)
	gamesCmd.Parse(os.Args[2:])

	conn, err := db.New()
	if err != nil {
		logger.Fatalf("failed to get db connection: %v", err)
	}

	games, err := conn.AllGames()
	if err != nil {
		logger.Fatalf("failed to show all games: %v", err)
	}
	logger.Printf("Past Games")
	for _, g := range games {
		logger.Printf("%s vs %s --- Score: %d - %d. Played %s", g.Player1Name, g.Player2Name, g.Player1Points, g.Player2Points, g.Played)
	}
}