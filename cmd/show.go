package cmd

import (
	"flag"
	"os"

	"github.com/dredly/ego/internal/db"
)

var showGames bool

func RunShow() {
	showCmd := flag.NewFlagSet("show", flag.ExitOnError)
	showCmd.BoolFunc("games", "shows a record of past games instead of the leaderboard", func(string) error {
		verboseLog("showing games instead of leaderboard")
		showGames = true
		return nil
	})
	addVerboseFlag(showCmd)
	showCmd.Parse(os.Args[2:])
	
	conn, err := db.New()
	if err != nil {
		logger.Fatalf("failed to get db connection: %v", err)
	}
	
	if showGames {
		showAllGames(conn)
	} else {
		showAllPlayers(conn)
	}
}

func showAllPlayers(conn *db.DBConnection) {
	players, err:= conn.Show()
	if err != nil {
		logger.Fatalf("failed to show all players: %v", err)
	}
	logger.Printf("Leaderboard")
	for i, player := range players {
		logger.Printf("%d. %s: %.2f", i + 1, player.Name, player.ELO)
	}
}

func showAllGames(conn *db.DBConnection) {

}