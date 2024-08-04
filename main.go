package main

import (
	"log"
	"os"

	"github.com/dredly/ego/cmd"
)

var logger = log.New(os.Stdout, "", 0)

func main() {
	if len(os.Args) < 2 {
		logger.Fatal("expected a subcommand")
    }

	switch os.Args[1] {
	case "create":
		cmd.RunCreate()
    case "add":
        cmd.RunAdd()
	case "leaderboard":
		cmd.RunLeaderboard()
	case "games":
		cmd.RunGames()
	case "record":
		cmd.RunRecord()
	case "undo":
		cmd.RunUndo()
	case "stats":
		cmd.RunStats()
    default:
		logger.Fatalf("unrecognised subcommand: %s\n", os.Args[1])
	}
}