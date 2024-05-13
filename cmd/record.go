package cmd

import (
	"flag"
	"os"

	"github.com/dredly/ego/internal/db"
)

func RunRecord() {
	if os.Args[2] == "draw" {
		recordDraw()
	} else {
		recordDecisive()
	}
}

func recordDecisive() {
	recordCmd := flag.NewFlagSet("record", flag.ExitOnError)
	addVerboseFlag(recordCmd)
	recordWinner := recordCmd.String("w", "", "name of the player who won")
	recordLoser := recordCmd.String("l", "", "name of the player who lost")
	multiplier := 1
	recordCmd.BoolFunc("donut", "whether the loser scored 0 points", func(string) error {
		verboseLog("recording a donut - setting multiplier to 2")
		multiplier = 2
		return nil
	})
	recordCmd.Parse(os.Args[2:])
	handleRecording(*recordWinner, *recordLoser, 1, multiplier)
}

func recordDraw() {
	recordDrawCmd := flag.NewFlagSet("draw", flag.ExitOnError)
	addVerboseFlag(recordDrawCmd)
	recordDrawCmd.Parse(os.Args[3:])
	if len(os.Args) < 4 {
		logger.Fatal("expected 2 player names when recording draw")
	}
	handleRecording(os.Args[3], os.Args[4], 0.5, 1)
}

func handleRecording(playerName1, playerName2 string, score float64, multiplier int) {
	conn, err := db.New()
	if err != nil {
		logger.Fatalf("failed to get db connection: %v", err)
	}

	player1, err := conn.FindPlayerByName(playerName1)
	if err != nil {
		logger.Fatalf("failed to find player by name: %v", err)
	}
	player2, err := conn.FindPlayerByName(playerName2)
	if err != nil {
		logger.Fatalf("failed to find player by name: %v", err)
	}

	player1ELOInitial := player1.ELO
	player2ELOInitial := player2.ELO

	player1.RecordResult(player2ELOInitial, score, multiplier)
	player2.RecordResult(player1ELOInitial,  1 - score, multiplier)

	if err := conn.UpdatePlayer(player1); err != nil {
		logger.Fatalf("failed to update elo: %v", err)
	}
	if err := conn.UpdatePlayer(player2); err != nil {
		logger.Fatalf("failed to update elo: %v", err)
	}

	if score == 0.5 {
		logger.Printf("recorded draw between %s and %s", playerName1, playerName2)
	} else {
		if multiplier == 2 {
			logger.Printf("recorded %s donut over %s\n", playerName1, playerName2)
		} else {
			logger.Printf("recorded %s win over %s\n", playerName1, playerName2)
		}
	}

	logger.Printf("%s elo: %.2f -> %.2f. %s elo: %.2f -> %.2f", playerName1, player1ELOInitial, player1.ELO, playerName2, player2ELOInitial, player2.ELO)
}