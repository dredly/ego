package main

import (
	"flag"
	"log"
	"os"
	"slices"

	"github.com/dredly/ego/internal/db"
)

var verbose bool;
var logger = log.New(os.Stdout, "", 0)

func main() {
	createCmd := flag.NewFlagSet("create", flag.ExitOnError)

	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	addName := addCmd.String("name", "", "name of the player to add to the leaderboard")

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

	switch args[1] {
	case "create":
		verboseLog("test message")
        createCmd.Parse(os.Args[2:])
		conn, err := db.New()
		if err != nil {
			log.Fatalf("failed to get db connection: %v", err)
		}
		err = conn.Initialise()
		if err != nil {
			logger.Fatalf("failed to initialise db: %v", err)
		}
		logger.Printf("Successfully initialised leaderboard")
    case "add":
        addCmd.Parse(os.Args[2:])
		logger.Printf("added new player %s", *addName)
	case "record":
		recordCmd.Parse(os.Args[2:])
		logger.Printf("recorded %s win over %s\n", *recordWinner, *recordLoser)
    default:
		logger.Fatalf("unrecognised subcommand: %s\n", os.Args[1])
	}
}

func verboseLog(msg string) {
	if verbose {
		logger.Print(msg)
	}
}