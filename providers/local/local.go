package local

import (
	"fmt"
	"log"
	"database/sql"
	_ "modernc.org/sqlite"
)

type Config struct {}

type LocalProvider struct {
	Config Config
	db *sql.DB
}

func (p *LocalProvider) Get(key string) (string, error) {
	var value string
	row := p.db.QueryRow(`SELECT value FROM secrets WHERE key = ?;`, key)
	err := row.Scan(&value)

	if err != nil {
		return "", err
	}

	fmt.Println("Secret fetched successfully!")
	return value, nil
}

func (p *LocalProvider) Set(key, value string) error {
	_, err := p.db.Exec(`
		INSERT INTO secrets (key, value)
		VALUES (?, ?)
		ON CONFLICT(key) DO UPDATE SET value = excluded.value;
	`, key, value)
	if err != nil {
		return err
	}

	fmt.Println("Secret added successfully!")
	return nil
}

func (p LocalProvider) Destroy(key string) error {
	result, err := p.db.Exec(`DELETE FROM secrets WHERE key = ?;`, key)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no record found for key: %s", key)
	}

	fmt.Println("Secret destroyed successfully!")
	return nil
}

func (p *LocalProvider) Setup() {
	db, err := sql.Open("sqlite", "tmp/example.db")
	if err != nil {
		log.Fatal(err)
	}
	p.db = db

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS secrets (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		key VARCHAR(64) UNIQUE NOT NULL,
		value VARCHAR(64) NOT NULL
	)`)
	if err != nil {
		log.Fatal(err)
	}
}

func NewProvider(c Config) *LocalProvider {
	p := LocalProvider{
		Config: c,
	}
	p.Setup()
	return &p
}
