package cmd

import (
	"flag"
	"os"

	"github.com/dredly/ego/internal/db"
	"github.com/dredly/ego/internal/types"
)

const defaultStartingELO = 1000

func RunAdd() {
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	name := addCmd.String("name", "", "name of the player to add to the leaderboard")
	elo := addCmd.Float64("elo", defaultStartingELO, "starting elo for the player")
	addVerboseFlag(addCmd)
	addCmd.Parse(os.Args[2:])
	if *elo <= 0 {
		logger.Fatal("ELO must be a positive number")
	}
	conn, err := db.New()
	if err != nil {
		logger.Fatalf("failed to get db connection: %v", err)
	}
	if err := conn.AddPlayer(*types.NewPlayer(*name, *elo)); err != nil {
		logger.Fatalf("failed to add player: %v", err)
	}
	logger.Printf("added new player %s with starting ELO %2.f", *name, *elo)
}