package cmd

import (
	"flag"
	"os"

	"github.com/dredly/ego/internal/db"
	"github.com/dredly/ego/internal/types"
)

func RunAdd() {
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	addName := addCmd.String("name", "", "name of the player to add to the leaderboard")
	addVerboseFlag(addCmd)
	addCmd.Parse(os.Args[2:])
	conn, err := db.New()
	if err != nil {
		logger.Fatalf("failed to get db connection: %v", err)
	}
	if err := conn.AddPlayer(*types.NewPlayer(*addName)); err != nil {
		logger.Fatalf("failed to add player: %v", err)
	}
	logger.Printf("added new player %s", *addName)
}