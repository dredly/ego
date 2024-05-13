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

	winner.RecordResult(loserELOInitial, 1, multiplier)
	loser.RecordResult(winnerELOInitial, 0, multiplier)

	if err := conn.UpdatePlayer(winner); err != nil {
		logger.Fatalf("failed to update winner elo: %v", err)
	}
	if err := conn.UpdatePlayer(loser); err != nil {
		logger.Fatalf("failed to update loser elo: %v", err)
	}

	if multiplier == 2 {
		logger.Printf("recorded %s donut over %s\n", *recordWinner, *recordLoser)
	} else {
		logger.Printf("recorded %s win over %s\n", *recordWinner, *recordLoser)
	}
	logger.Printf("%s elo: %.2f -> %.2f. %s elo: %.2f -> %.2f", *recordWinner, winnerELOInitial, winner.ELO, *recordLoser, loserELOInitial, loser.ELO)
}

func recordDraw() {
	recordDrawCmd := flag.NewFlagSet("draw", flag.ExitOnError)
	addVerboseFlag(recordDrawCmd)
	recordDrawCmd.Parse(os.Args[3:])
	if len(os.Args) < 4 {
		logger.Fatal("expected 2 player names when recording draw")
	}
	conn, err := db.New()
	if err != nil {
		logger.Fatalf("failed to get db connection: %v", err)
	}
	player1, err := conn.FindPlayerByName(os.Args[3])
	if err != nil {
		logger.Fatalf("failed to find player by name: %v", err)
	}
	player2, err := conn.FindPlayerByName(os.Args[4])
	if err != nil {
		logger.Fatalf("failed to find player by name: %v", err)
	}

	player1ELOInitial := player1.ELO
	player2ELOInitial := player2.ELO

	player1.RecordResult(player2ELOInitial, 0.5, 1)
	player2.RecordResult(player1ELOInitial, 0.5, 1)

	if err := conn.UpdatePlayer(player1); err != nil {
		logger.Fatalf("failed to update elo: %v", err)
	}
	if err := conn.UpdatePlayer(player2); err != nil {
		logger.Fatalf("failed to update elo: %v", err)
	}

	logger.Printf("recorded draw between %s and %s", player1.Name, player2.Name)
	logger.Printf("%s elo: %.2f -> %.2f. %s elo: %.2f -> %.2f", player1.Name, player1ELOInitial, player1.ELO, player2.Name, player2ELOInitial, player2.ELO)
}