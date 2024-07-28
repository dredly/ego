package cmd

import (
	"flag"
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

	peakELO, err := conn.PeakELOForPlayer(*name)
	if err != nil {
		logger.Fatalf("failed to determine peak ELO for player: %v", err)
	}

	logger.Printf("Stats for %s", *name)
	logger.Printf("Peak ELO: %.4f", peakELO)
}