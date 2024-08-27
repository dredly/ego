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

func (conn DBConnection) Games(playerName string, limit uint) ([]types.GameSummary, error) {
	all_games_query := `
		SELECT player1name, player1points, player2name, player2points, played
		FROM game_summaries
	`
	ordering := " ORDER BY played DESC"

	var playerFilter string
	var limitClause string
	args := []any{}

	if playerName != "" {
		args = append(args, playerName)
		playerFilter = fmt.Sprintf(" WHERE player1name = $%d OR player2name = $%d", len(args), len(args)) 
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

	var gameSummaries []types.GameSummary
	for rows.Next() {
		var gs types.GameSummary
		rows.Scan(&gs.Player1Name, &gs.Player1Points, &gs.Player2Name, &gs.Player2Points, &gs.Played)
		gameSummaries = append(gameSummaries, gs)	
	}
	return gameSummaries, nil
}

func (conn DBConnection) UndoGame() (types.GameDetail, error) {
	tx, err := conn.db.Begin()
    if err != nil {
        return types.GameDetail{}, err
    }

	latestGameSQL := `
		SELECT
			id,
			player1name, player1elobefore, player1eloafter,
			player2name, player2elobefore, player2eloafter,
			played
		FROM game_details
		ORDER BY played DESC
		LIMIT 1
	`
	conn.logSQL(latestGameSQL)
	row := conn.db.QueryRow(latestGameSQL)
	var gd types.GameDetail
	err = row.Scan(
		&gd.ID,
		&gd.Player1Name, &gd.Player1ELOBefore, &gd.Player1ELOAfter,
		&gd.Player2Name, &gd.Player2ELOBefore, &gd.Player2ELOAfter,
		&gd.Played,
	)
	if err != nil {
		tx.Rollback()
		return types.GameDetail{}, err
	}

	updatePlayerSQL := `UPDATE players SET elo = $1 WHERE name = $2`
	conn.logSQL(updatePlayerSQL)
	_, err = tx.Exec(updatePlayerSQL, gd.Player1ELOAfter, gd.Player1Name)
	if err != nil {
		tx.Rollback()
		return types.GameDetail{}, err
	}
	_, err = tx.Exec(updatePlayerSQL, gd.Player2ELOAfter, gd.Player2Name)
	if err != nil {
		tx.Rollback()
		return types.GameDetail{}, err
	}

	deletePlayerGamesSQL := `DELETE FROM player_games WHERE gameid = $1`
	conn.logSQL(deletePlayerGamesSQL)
	_, err = tx.Exec(deletePlayerGamesSQL, gd.ID)
	if err != nil {
		tx.Rollback()
		return types.GameDetail{}, err
	}

	deleteGameSQL := `DELETE FROM games WHERE id = $1`
	conn.logSQL(deleteGameSQL)
	_, err = tx.Exec(deleteGameSQL, gd.ID)
	if err != nil {
		tx.Rollback()
		return types.GameDetail{}, err
	}

	err = tx.Commit()
    if err != nil {
        return types.GameDetail{}, err
    }

	return gd, err
}