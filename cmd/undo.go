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

	deletedGame, err := conn.DeleteMostRecentGame()
	if err != nil {
		logger.Fatalf("failed to delete game: %v", err)
	}

	player1Reverted, err := conn.UpdateELOByID(deletedGame.Player1ID, deletedGame.Player1ELOBefore)
	if err != nil {
		logger.Fatalf("failed to update ELO for player: %v", err)
	}
	player2Reverted, err := conn.UpdateELOByID(deletedGame.Player2ID, deletedGame.Player2ELOBefore)
	if err != nil {
		logger.Fatalf("failed to update ELO for player: %v", err)
	}

	logger.Printf("Reverted last game between %s and %s, played on %s", player1Reverted.Name, player2Reverted.Name, deletedGame.Played)
	logger.Printf("%s elo: %.2f -> %.2f. %s elo: %.2f -> %.2f", 
		player1Reverted.Name, deletedGame.Player1ELOAfter, deletedGame.Player1ELOBefore, 
		player2Reverted.Name, deletedGame.Player2ELOAfter, deletedGame.Player2ELOBefore,
	)
}