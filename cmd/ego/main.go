package main

import (
	"flag"
	"log"
	"os"
	"slices"
	"strings"

	"github.com/dredly/ego/internal/db"
	"github.com/dredly/ego/internal/types"
)

var verbose bool;
var logger = log.New(os.Stdout, "", 0)

func main() {
	createCmd := flag.NewFlagSet("create", flag.ExitOnError)

	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	addName := addCmd.String("name", "", "name of the player to add to the leaderboard")

	showCmd := flag.NewFlagSet("show", flag.ExitOnError)

	recordCmd := flag.NewFlagSet("record", flag.ExitOnError)
	recordWinner := recordCmd.String("w", "", "name of the player who won")
	recordLoser := recordCmd.String("l", "", "name of the player who lost")

	if len(os.Args) < 2 {
		logger.Fatal("expected a subcommand")
    }

	var verboseArgIdx int;
	for i, arg := range os.Args {
		if arg == "-v" || arg == "verbose" {
			verbose = true;
			verboseArgIdx = i;
			break;
		}
	}
	var args []string
	if verbose {
		args = slices.Delete(os.Args, verboseArgIdx, verboseArgIdx + 1)
	} else {
		args = os.Args
	}

	verboseLog("args given: " + strings.Join(args, ", "))

	switch args[1] {
	case "create":
        createCmd.Parse(args[2:])
		conn, err := db.New()
		if err != nil {
			logger.Fatalf("failed to get db connection: %v", err)
		}
		if err := conn.Initialise(); err != nil {
			logger.Fatalf("failed to initialise db: %v", err)
		}
		logger.Printf("Successfully initialised leaderboard")
    case "add":
        addCmd.Parse(args[2:])
		conn, err := db.New()
		if err != nil {
			logger.Fatalf("failed to get db connection: %v", err)
		}
		if err := conn.AddPlayer(*types.NewPlayer(*addName)); err != nil {
			logger.Fatalf("failed to add player: %v", err)
		}
		logger.Printf("added new player %s", *addName)
	case "show":
		showCmd.Parse(args[2:])
		conn, err := db.New()
		if err != nil {
			logger.Fatalf("failed to get db connection: %v", err)
		}
		players, err:= conn.Show()
		if err != nil {
			logger.Fatalf("failed to show all players: %v", err)
		}
		logger.Printf("Leaderboard")
		for i, player := range players {
			logger.Printf("%d. %s: %.2f", i + 1, player.Name, player.ELO)
		}
	case "record":
		recordCmd.Parse(args[2:])
		conn, err := db.New()
		if err != nil {
			logger.Fatalf("failed to get db connection: %v", err)
		}
		winner, err := conn.FindPlayerByName(*recordWinner)
		if err != nil {
			logger.Fatalf("failed to find winner by name: %v", err)
		}
		loser, err := conn.FindPlayerByName(*recordLoser)
		if err != nil {
			logger.Fatalf("failed to find loser by name: %v", err)
		}

		winnerELOInitial := winner.ELO
		loserELOInitial := loser.ELO

		winner.RecordWin(loserELOInitial)
		loser.RecordLoss(winnerELOInitial)

		if err := conn.UpdatePlayer(winner); err != nil {
			logger.Fatalf("failed to update winner elo: %v", err)
		}
		if err := conn.UpdatePlayer(loser); err != nil {
			logger.Fatalf("failed to update loser elo: %v", err)
		}

		logger.Printf("recorded %s win over %s\n", *recordWinner, *recordLoser)
		logger.Printf("%s elo: %.2f -> %.2f. %s elo: %.2f -> %.2f", *recordWinner, winnerELOInitial, winner.ELO, *recordLoser, loserELOInitial, loser.ELO)
    default:
		logger.Fatalf("unrecognised subcommand: %s\n", args[1])
	}
}

func verboseLog(msg string) {
	if verbose {
		logger.Print(msg)
	}
}