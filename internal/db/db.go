package db

import (
	"database/sql"
	"os"
	"path/filepath"

	"github.com/dredly/ego/internal/types"
	_ "github.com/mattn/go-sqlite3"
)

type DBConnection struct {
	db *sql.DB
}

func New() (*DBConnection, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	dirExists, err := exists(filepath.Join(home, ".ego"))
	if err != nil {
		return nil, err
	}
	if !dirExists {
		err := os.MkdirAll(filepath.Join(home, ".ego"), os.ModePerm)
		if err != nil {
			return nil, err
		}
	}
	pathToDB := filepath.Join(home, ".ego", "ego.db")
	dbFileExists, err := exists(pathToDB)
	if err != nil {
		return nil, err
	}
	if !dbFileExists {
		f, err := os.Create(pathToDB)
		if err != nil {
			return nil, err
		}
		defer f.Close()
	}
	db, err := sql.Open("sqlite3", pathToDB)
	if err != nil {
		return nil, err
	}
	return &DBConnection{db}, nil
}

func (conn DBConnection) Initialise() error {
	_, err := conn.db.Exec(`CREATE TABLE IF NOT EXISTS players (
		id    INTEGER PRIMARY KEY,
        name  TEXT NOT NULL,
		elo   REAL,
		CONSTRAINT name_unique UNIQUE (name)
	);
	CREATE TABLE IF NOT EXISTS games (
		id INTEGER PRIMARY KEY,
		player1id INTEGER,
		player2id INTEGER,
		player1points INTEGER,
		player2points INTEGER,
		player1elobefore REAL,
		player2elobefore REAL,
		player1eloafter REAL,
		player2eloafter REAL,
		played DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(player1id) REFERENCES players(id) ON DELETE SET NULL,
		FOREIGN KEY(player2id) REFERENCES players(id) ON DELETE SET NULL
	);`)
	if err != nil {
		return err
	}
	return nil
}

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

func (conn DBConnection) Show() ([]types.Player, error) {
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

func exists(path string) (bool, error) {
    _, err := os.Stat(path)
    if err == nil { return true, nil }
    if os.IsNotExist(err) { return false, nil }
    return false, err
}