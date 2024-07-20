package cmd

import (
	"flag"
	"os"

	"github.com/dredly/ego/internal/db"
)

func RunUndo() {
	undoCmd := flag.NewFlagSet("undo", flag.ExitOnError)
	addVerboseFlag(undoCmd)
	addDbPathFlag(undoCmd)
	undoCmd.Parse(os.Args[2:])

	conn, err := db.Connect(dbPath, verbose)
	if err != nil {
		logger.Fatalf("failed to get db connection: %v", err)
	}

	mostRecentGame, err := conn.MostRecentGame()
	_ = mostRecentGame
}