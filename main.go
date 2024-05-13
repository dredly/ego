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
	case "show":
		cmd.RunShow()
	case "record":
		cmd.RunRecord()
    default:
		logger.Fatalf("unrecognised subcommand: %s\n", os.Args[1])
	}
}