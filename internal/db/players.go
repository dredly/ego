package db

import "github.com/dredly/ego/internal/types"

func (conn DBConnection) AddPlayer(p types.Player) error {
	sql := "INSERT INTO players (name, elo) values ($1, $2)"
	conn.logSQL(sql)
	stmt, err := conn.db.Prepare(sql)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(p.Name, p.ELO)
	if err != nil {
		return err
	}
	return nil
}

func (conn DBConnection) AllPlayers() ([]types.Player, error) {
	sql := "SELECT id, name, elo FROM players ORDER BY elo DESC"
	conn.logSQL(sql)
	rows, err := conn.db.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var players []types.Player
	for rows.Next() {
		var p types.Player
		rows.Scan(&p.ID, &p.Name, &p.ELO)
		players = append(players, p)
	}
	return players, nil
}

func (conn DBConnection) FindPlayerByName(name string) (types.Player, error) {
	sql := "SELECT id, elo FROM players WHERE name = $1"
	conn.logSQL(sql)
	row := conn.db.QueryRow(sql, name)
	p := types.Player{
		Name: name,
	}
	err := row.Scan(&p.ID, &p.ELO)
	if err != nil {
		return types.Player{}, err
	}
	return p, nil
}

func (conn DBConnection) UpdatePlayer(p types.Player) error {
	sql := "UPDATE players SET elo = $1 WHERE name = $2"
	conn.logSQL(sql)
	stmt, err := conn.db.Prepare(sql)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(p.ELO, p.Name)
	if err != nil {
		return err
	}
	return nil
}

func (conn DBConnection) PeakELOForPlayer(name string) (float64, error) {
	sql := `
		SELECT MAX(MAX(pg.elobefore, pg.eloafter)) AS peak_elo
		FROM player_games AS pg
		INNER JOIN players AS p ON pg.playerid = p.id
		WHERE p.name = $1
	`
	conn.logSQL(sql)

	row := conn.db.QueryRow(sql, name)
	var peakELO float64

	err := row.Scan(&peakELO)
	if err != nil {
		return 0, err
	}

	return peakELO, nil
}

func (conn DBConnection) GameResults(name string) (types.GameResults, error) {
	sql := `
		SELECT 
			SUM(CASE 
				WHEN rpg.points > opponent.points THEN 1 
				ELSE 0 
			END) AS won,
			SUM(CASE 
				WHEN rpg.points = opponent.points THEN 1 
				ELSE 0 
			END) AS drawn,
			SUM(CASE 
				WHEN rpg.points < opponent.points THEN 1 
				ELSE 0 
			END) AS lost
		FROM ranked_player_games rpg
		JOIN players p ON rpg.playerid = p.id
		JOIN ranked_player_games opponent ON rpg.gameid = opponent.gameid AND rpg.playerid != opponent.playerid
		WHERE p.name = $1;
	`
	conn.logSQL(sql)

	row := conn.db.QueryRow(sql, name)
	var gr types.GameResults
	err := row.Scan(&gr.Won, &gr.Drawn, &gr.Lost)
	if err != nil {
		return types.GameResults{}, err
	}

	return gr, nil
}