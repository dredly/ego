package cmd

import (
	"flag"
	"os"

	"github.com/dredly/ego/internal/db"
)

func RunCreate() {
	createCmd := flag.NewFlagSet("create", flag.ExitOnError)
	addVerboseFlag(createCmd)
	addDbPathFlag(createCmd)
	createCmd.Parse(os.Args[2:])
	conn, err := db.New(dbPath, verbose)
	if err != nil {
		logger.Fatalf("failed to get db connection: %v", err)
	}
	if err := conn.Initialise(); err != nil {
		logger.Fatalf("failed to initialise db: %v", err)
	}
	logger.Printf("Successfully initialised leaderboard")
}