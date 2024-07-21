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

func (conn DBConnection) UpdateELOByID(id int, newELO float64) (types.Player, error) {
	sql := `
		UPDATE players SET elo = $1
		WHERE id = $2
		RETURNING name, elo
	`
	conn.logSQL(sql)

	row := conn.db.QueryRow(sql, newELO, id)
	p := types.Player{
		ID: id,
	}
	err := row.Scan(&p.Name, &p.ELO)
	if err != nil {
		return types.Player{}, err
	}
	return p, nil
}