package cmd

import (
	"flag"
	"os"
	"strconv"
	"strings"

	"github.com/dredly/ego/internal/db"
	"github.com/dredly/ego/internal/types"
)

func RunRecord() {
	recordCmd := flag.NewFlagSet("record", flag.ExitOnError)
	addVerboseFlag(recordCmd)
	addDbPathFlag(recordCmd)
	playerName1 := recordCmd.String("p1", "", "name of player 1")
	playerName2 := recordCmd.String("p2", "", "name of player 2")
	score := recordCmd.String("score", "", "final score of the game, 2 integers separated by a dash, e.g. '3-2'")
	recordCmd.Parse(os.Args[2:])
	scores := strings.Split(*score, "-")
	if len(scores) != 2 {
		logger.Fatal("invalid format for score, should be 2 integers separated by a dash")
	}
	p1Points, err := strconv.Atoi(scores[0])
	if err != nil {
		logger.Fatalf("failed to parse points for player 1: %v", err)
	}
	p2Points, err := strconv.Atoi(scores[1])
	if err != nil {
		logger.Fatalf("failed to parse points for player 2: %v", err)
	}

	conn, err := db.Connect(dbPath, verbose)
	if err != nil {
		logger.Fatalf("failed to get db connection: %v", err)
	}

	player1, err := conn.FindPlayerByName(*playerName1)
	if err != nil {
		logger.Fatalf("failed to find player by name: %v", err)
	}
	player2, err := conn.FindPlayerByName(*playerName2)
	if err != nil {
		logger.Fatalf("failed to find player by name: %v", err)
	}

	player1ELOInitial := player1.ELO
	player2ELOInitial := player2.ELO

	eloScore := eloScore(p1Points, p2Points)
	multiplier := multiplier(p1Points, p2Points)

	player1.RecordResult(player2ELOInitial, eloScore, multiplier)
	player2.RecordResult(player1ELOInitial, 1 - eloScore, multiplier)

	game := types.Game{
		Player1ID: player1.ID,
		Player2ID: player2.ID,
		Player1Points: p1Points,
		Player2Points: p2Points,
		Player1ELOBefore: player1ELOInitial,
		Player2ELOBefore: player2ELOInitial,
		Player1ELOAfter: player1.ELO,
		Player2ELOAfter: player2.ELO,
	}
	
	if err := conn.AddGame(game); err != nil {
		logger.Fatalf("failed to add game: %v", err)
	}
	if err := conn.UpdatePlayer(player1); err != nil {
		logger.Fatalf("failed to update elo: %v", err)
	}
	if err := conn.UpdatePlayer(player2); err != nil {
		logger.Fatalf("failed to update elo: %v", err)
	}

	logger.Printf("Recorded result %s between %s and %s", *score, *playerName1, *playerName2)
	logger.Printf("%s elo: %.2f -> %.2f. %s elo: %.2f -> %.2f", *playerName1, player1ELOInitial, player1.ELO, *playerName2, player2ELOInitial, player2.ELO)
}

func eloScore(p1Points, p2Points int) float64 {
	if p1Points > p2Points {
		return 1
	}
	if p1Points < p2Points {
		return 0
	}
	return 0.5
}

func multiplier(p1Points, p2Points int) int {
	if p1Points == 0 || p2Points == 0 {
		return 2
	}
	return 1
}