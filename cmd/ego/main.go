package main

import (
	"flag"
	"log"
	"os"

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
	var donut bool
	recordCmd.BoolFunc("donut", "whether the loser scored 0 points", func(string) error {
		donut = true
		return nil
	})

	subCommands := []*flag.FlagSet{createCmd, addCmd, showCmd, recordCmd}
	for _, sc := range subCommands {
		sc.BoolFunc("verbose", "enable more detailed logging", func(string) error {
			verbose = true
			return nil
		})
	}

	if len(os.Args) < 2 {
		logger.Fatal("expected a subcommand")
    }

	switch os.Args[1] {
	case "create":
        createCmd.Parse(os.Args[2:])
		conn, err := db.New()
		if err != nil {
			logger.Fatalf("failed to get db connection: %v", err)
		}
		if err := conn.Initialise(); err != nil {
			logger.Fatalf("failed to initialise db: %v", err)
		}
		logger.Printf("Successfully initialised leaderboard")
    case "add":
        addCmd.Parse(os.Args[2:])
		conn, err := db.New()
		if err != nil {
			logger.Fatalf("failed to get db connection: %v", err)
		}
		if err := conn.AddPlayer(*types.NewPlayer(*addName)); err != nil {
			logger.Fatalf("failed to add player: %v", err)
		}
		logger.Printf("added new player %s", *addName)
	case "show":
		showCmd.Parse(os.Args[2:])
		verboseLog("verbose log in show command")
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
		recordCmd.Parse(os.Args[2:])
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

		winner.RecordWin(loserELOInitial, donut)
		loser.RecordLoss(winnerELOInitial, donut)

		if err := conn.UpdatePlayer(winner); err != nil {
			logger.Fatalf("failed to update winner elo: %v", err)
		}
		if err := conn.UpdatePlayer(loser); err != nil {
			logger.Fatalf("failed to update loser elo: %v", err)
		}

		if donut {
			logger.Printf("recorded %s donut over %s\n", *recordWinner, *recordLoser)
		} else {
			logger.Printf("recorded %s win over %s\n", *recordWinner, *recordLoser)
		}
		logger.Printf("%s elo: %.2f -> %.2f. %s elo: %.2f -> %.2f", *recordWinner, winnerELOInitial, winner.ELO, *recordLoser, loserELOInitial, loser.ELO)
    default:
		logger.Fatalf("unrecognised subcommand: %s\n", os.Args[1])
	}
}

func verboseLog(msg string) {
	if verbose {
		logger.Print(msg)
	}
}