package db

import (
	"fmt"
	"time"

	"github.com/dredly/ego/internal/types"
)

func (conn DBConnection) AddGame(g types.Game) error {
	stmt, err := conn.db.Prepare(`INSERT INTO games (
		player1id,
		player2id,
		player1points,
		player2points,
		player1elobefore,
		player2elobefore,
		player1eloafter,
		player2eloafter
	) values ($1, $2, $3, $4, $5, $6, $7, $8)`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(g.Player1ID, g.Player2ID, g.Player1Points, g.Player2Points, g.Player1ELOBefore, g.Player2ELOBefore, g.Player1ELOAfter, g.Player2ELOAfter) 
	if err != nil {
		return err
	}
	return nil
}

func (conn DBConnection) Games(playerName string, limit uint) ([]types.GameDisplay, error) {
	all_games_query := `
		SELECT p1.name, p2.name, g.player1Points, g.player2Points, g.played 
		FROM games AS g 
		INNER JOIN players as p1 ON g.player1id = p1.id
		INNER JOIN players as p2 ON g.player2id = p2.id
	`
	ordering := " ORDER BY g.played DESC"

	var playerFilter string
	var limitClause string
	args := []any{}

	if playerName != "" {
		args = append(args, playerName)
		playerFilter = fmt.Sprintf(" WHERE p1.name = $%d OR p2.name = $%d", len(args), len(args)) 
	}
	if limit > 0 {
		args = append(args, limit)
		limitClause = fmt.Sprintf(" LIMIT $%d", len(args)) 
	}

	rows, err := conn.db.Query(all_games_query + playerFilter + ordering + limitClause, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var games []types.GameDisplay
	for rows.Next() {
		var p1Name, p2Name string
		var p1Points, p2Points int
		var played time.Time
		rows.Scan(&p1Name, &p2Name, &p1Points, &p2Points, &played)
		games = append(games, types.GameDisplay{
			Player1Name: p1Name,
			Player2Name: p2Name,
			Player1Points: p1Points,
			Player2Points: p2Points,
			Played: played,
		})	
	}
	return games, nil
}