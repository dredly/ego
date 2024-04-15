package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/dredly/ego/internal/db"
)

func main() {
	createCmd := flag.NewFlagSet("create", flag.ExitOnError)

	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	addName := addCmd.String("name", "", "name of the player to add to the leaderboard")

	recordCmd := flag.NewFlagSet("record", flag.ExitOnError)
	recordWinner := recordCmd.String("w", "", "name of the player who won")
	recordLoser := recordCmd.String("l", "", "name of the player who lost")

	if len(os.Args) < 2 {
        fmt.Println("expected a subcommand")
        os.Exit(1)
    }

	switch os.Args[1] {
	case "create":
        createCmd.Parse(os.Args[2:])
		conn, err := db.New()
		if err != nil {
			fmt.Printf("Failed to get db connection: %v\n", err)
			os.Exit(1)
		}
		err = conn.Initialise()
		if err != nil {
			fmt.Printf("Failed to initialise database: %v/n", err)
			os.Exit(1)
		}
		fmt.Println("Initialised database")
    case "add":
        addCmd.Parse(os.Args[2:])
		fmt.Printf("adding new player %s\n", *addName)
	case "record":
		recordCmd.Parse(os.Args[2:])
		fmt.Printf("recording %s win over %s\n", *recordWinner, *recordLoser)
    default:
        fmt.Printf("unrecognised subcommand: %s\n", os.Args[1])
        os.Exit(1)
	}
}