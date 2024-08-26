package db

import (
	"fmt"

	"github.com/dredly/ego/internal/types"
)

func (conn DBConnection) RecordGame(gr types.GameRecording) error {
	tx, err := conn.db.Begin()
    if err != nil {
        return err
    }

	updatePlayerSQL := "UPDATE players SET elo = $1 WHERE name = $2"
	conn.logSQL(updatePlayerSQL)
	_, err = tx.Exec(updatePlayerSQL, gr.Player1.ELOAfter, gr.Player1.Player.Name)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.Exec(updatePlayerSQL, gr.Player2.ELOAfter, gr.Player2.Player.Name)
	if err != nil {
		tx.Rollback()
		return err
	}

	insertGameSQL := "INSERT INTO games DEFAULT VALUES"
	conn.logSQL(insertGameSQL)
    result, err := tx.Exec(insertGameSQL)
    if err != nil {
        tx.Rollback()
        return err
    }
    lastGameID, err := result.LastInsertId()
    if err != nil {
        tx.Rollback()
        return err
    }

	insertPlayerGameSQL := `
		INSERT INTO player_games
		(gameid, playerid, points, elobefore, eloafter)
		VALUES ($1, $2, $3, $4, $5), ($1, $6, $7, $8, $9)
	`
	conn.logSQL(insertPlayerGameSQL)
	_, err = tx.Exec(insertPlayerGameSQL, 
		lastGameID,
		gr.Player1.Player.ID, gr.Player1.Points, gr.Player1.ELOBefore, gr.Player1.ELOAfter,
		gr.Player2.Player.ID, gr.Player2.Points, gr.Player2.ELOBefore, gr.Player2.ELOAfter,
	)
	if err != nil {
        tx.Rollback()
        return err
    }

	err = tx.Commit()
    if err != nil {
        return err
    }

	return nil
}

func (conn DBConnection) Games(playerName string, limit uint) ([]types.GameDisplay, error) {
	all_games_query := `
		SELECT g.played 
		FROM games AS g 
		INNER JOIN player_games as pg ON g.id = pg.gameid

		INNER JOIN players as p ON g.playerid = p.id
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

func (conn DBConnection) DeleteMostRecentGame() (types.Game, error) {
	sql := `
		DELETE
		FROM games
		WHERE played = (SELECT MAX(played) FROM games)
		RETURNING
			id,
			player1id, player2id, 
			player1points, player2points, 
			player1elobefore, player2elobefore,  
			player1eloafter, player2eloafter,
			played
	`
	conn.logSQL(sql)

	row := conn.db.QueryRow(sql)
	var g types.Game
	err := row.Scan(
		&g.ID, 
		&g.Player1ID, &g.Player2ID, 
		&g.Player1Points, &g.Player2Points, 
		&g.Player1ELOBefore, &g.Player2ELOBefore, 
		&g.Player1ELOAfter, &g.Player2ELOAfter, 
		&g.Played,
	)

	if err != nil {
		return types.Game{}, err
	}
	return g, nil
}