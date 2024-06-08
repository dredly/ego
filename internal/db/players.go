package db

import "github.com/dredly/ego/internal/types"

func (conn DBConnection) AddPlayer(p types.Player) error {
	stmt, err := conn.db.Prepare("INSERT INTO players (name, elo) values ($1, $2)")
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
	rows, err := conn.db.Query("SELECT name, elo FROM players ORDER BY elo DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var players []types.Player
	for rows.Next() {
		var name string
		var elo float64
		rows.Scan(&name, &elo)
		players = append(players, types.Player{Name: name, ELO: elo})
	}
	return players, nil
}

func (conn DBConnection) FindPlayerByName(name string) (types.Player, error) {
	row := conn.db.QueryRow("SELECT id, elo FROM players WHERE name = $1", name)
	var id int
	var elo float64
	err := row.Scan(&id, &elo)
	if err != nil {
		return types.Player{}, err
	}
	return types.Player{ID: id, Name: name, ELO: elo}, nil
}

func (conn DBConnection) UpdatePlayer(p types.Player) error {
	stmt, err := conn.db.Prepare("UPDATE players SET elo = $1 WHERE name = $2")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(p.ELO, p.Name)
	if err != nil {
		return err
	}
	return nil
}