package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/dredly/ego/internal/db"
)

func RunGames() {
	gamesCmd := flag.NewFlagSet("games", flag.ExitOnError)
	limit := gamesCmd.Uint("limit", 0, "number of games to show, will show all by default or if set to 0")
	player := gamesCmd.String("player", "", "if specified, will only show games for that player")
	addVerboseFlag(gamesCmd)
	gamesCmd.Parse(os.Args[2:])

	conn, err := db.New()
	if err != nil {
		logger.Fatalf("failed to get db connection: %v", err)
	}

	games, err := conn.Games(*player, *limit)
	if err != nil {
		logger.Fatalf("failed to show all games: %v", err)
	}
	logger.Printf(header(*player))
	for _, g := range games {
		logger.Printf("%s vs %s --- Score: %d - %d. Played %s", g.Player1Name, g.Player2Name, g.Player1Points, g.Player2Points, g.Played)
	}
}

func header(playerName string) string {
	if playerName == "" {
		return "Past Games"
	}
	return fmt.Sprintf("%s's past games", playerName)
}