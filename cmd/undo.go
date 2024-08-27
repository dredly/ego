package cmd

import (
	"flag"
	"os"

	"github.com/dredly/ego/internal/db"
)

func RunUndo() {
	undoCmd := flag.NewFlagSet("undo", flag.ExitOnError)
	addVerboseFlag(undoCmd)
	addDbPathFlag(undoCmd)
	undoCmd.Parse(os.Args[2:])

	conn, err := db.Connect(dbPath, verbose)
	if err != nil {
		logger.Fatalf("failed to get db connection: %v", err)
	}

	deletedGameDetail, err := conn.UndoGame()
	if err != nil {
		logger.Fatalf("failed to undo game: %v", err)
	}

	logger.Printf("Reverted last game between %s and %s, played on %s", 
		deletedGameDetail.Player1Name, deletedGameDetail.Player2Name, deletedGameDetail.Played,
	)
	logger.Printf("%s elo: %.2f -> %.2f. %s elo: %.2f -> %.2f", 
		deletedGameDetail.Player1Name, deletedGameDetail.Player1ELOAfter, deletedGameDetail.Player1ELOBefore, 
		deletedGameDetail.Player2Name, deletedGameDetail.Player2ELOAfter, deletedGameDetail.Player2ELOBefore,
	)
}