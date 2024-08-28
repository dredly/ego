package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/dredly/ego/internal/db"
)

var (
	statsSpecified bool
	played bool
	winRate bool
	peak bool
)

func RunStats() {
	statsCmd := flag.NewFlagSet("stats", flag.ExitOnError)
	name := statsCmd.String("name", "", "name of the player to display stats for")
	statsCmd.BoolFunc("played", "show number of games played", func(string) error {
		statsSpecified = true
		played = true
		return nil
	})
	statsCmd.BoolFunc("winrate", "show win rate", func(string) error {
		statsSpecified = true
		winRate = true
		return nil
	})
	statsCmd.BoolFunc("peak", "show peak ELO", func(string) error {
		statsSpecified = true
		peak = true
		return nil
	})
	addVerboseFlag(statsCmd)
	addDbPathFlag(statsCmd)
	statsCmd.Parse(os.Args[2:])

	if statsSpecified {
		verboseLog("stats have been specified")
	} else {
		verboseLog("displaying all stats")
	}

	conn, err := db.Connect(dbPath, verbose)
	if err != nil {
		logger.Fatalf("failed to get db connection: %v", err)
	}

	gameResults, err := conn.GameResults(*name)
	if err != nil {
		logger.Fatalf("Failed to query game results for player: %v", err)
	}

	outputLines := []string{
		fmt.Sprintf("Stats for %s", *name),
	}

	if played || !statsSpecified {
		outputLines = append(outputLines, fmt.Sprintf("Games Played: %d", gameResults.Total()))
	}

	if winRate || !statsSpecified {
		outputLines = append(outputLines, fmt.Sprintf("Win rate: %s", gameResults.WinRateStr()))
	}

	if peak || !statsSpecified {
		peakELO, err := conn.PeakELOForPlayer(*name)
		if err != nil {
			logger.Fatalf("failed to determine peak ELO for player: %v", err)
		}
		outputLines = append(outputLines, fmt.Sprintf("Peak ELO: %.4f", peakELO))
	}

	for _, line := range outputLines {
		logger.Print(line)
	}
}