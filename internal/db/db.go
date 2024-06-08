package db

import (
	"database/sql"
	"os"
	"path/filepath"

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

func exists(path string) (bool, error) {
    _, err := os.Stat(path)
    if err == nil { return true, nil }
    if os.IsNotExist(err) { return false, nil }
    return false, err
}