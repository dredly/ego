package db

import (
	"fmt"

	"github.com/dredly/ego/internal/types"
)

func (conn DBConnection) AddGame(g types.Game) error {
	sql := `INSERT INTO games (
		player1id,
		player2id,
		player1points,
		player2points,
		player1elobefore,
		player2elobefore,
		player1eloafter,
		player2eloafter
	) values ($1, $2, $3, $4, $5, $6, $7, $8)`
	conn.logSQL(sql)
	stmt, err := conn.db.Prepare(sql)
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

	sql := all_games_query + playerFilter + ordering + limitClause
	conn.logSQL(sql)

	rows, err := conn.db.Query(sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var games []types.GameDisplay
	for rows.Next() {
		var g types.GameDisplay
		rows.Scan(&g.Player1Name, &g.Player2Name, &g.Player1Points, &g.Player2Points, &g.Played)
		games = append(games, g)	
	}
	return games, nil
}

func (conn DBConnection) MostRecentGame() (types.Game, error) {
	sql := `
		SELECT id, player1id, player2id, player1points, player2points, player1name, player2name, player1elobefore, player2elobefore 
		FROM games 
		ORDER BY played DESC 
		LIMIT 1
	`
	conn.logSQL(sql)

	// TODO: complete this
	return types.Game{}, nil
}