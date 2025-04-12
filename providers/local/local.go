package local

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite" // Driver necessary but never directly called
)

type Config struct {
	// TODO (ngeorge): Add options for provider configuration
	SqlitePath string
}

type Provider struct {
	ProviderConfig Config
	db             *sql.DB
}

func (p *Provider) Get(key string) (string, error) {
	var value string
	row := p.db.QueryRow(`SELECT value FROM secrets WHERE key = ?;`, key)
	err := row.Scan(&value)

	if err != nil {
		return "", fmt.Errorf("unable to get secret with key \"%s\" %w", key, err)
	}

	return value, nil
}

func (p *Provider) Exists(key string) (bool, error) {
	var value string
	row := p.db.QueryRow(`SELECT id FROM secrets WHERE key = ?;`, key)
	err := row.Scan(&value)

	if err == sql.ErrNoRows {
		return false, nil
	}

	if err != nil {
		return false, fmt.Errorf("error checking if secret exists with key \"%s\": %w", key, err)
	}

	return true, nil
}

func (p *Provider) Set(key, value string) error {
	_, err := p.db.Exec(`
		INSERT INTO secrets (key, value)
		VALUES (?, ?)
		ON CONFLICT(key) DO UPDATE SET value = excluded.value;
	`, key, value)
	if err != nil {
		return err
	}

	return nil
}

func (p *Provider) Destroy(key string) error {
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

	return nil
}

func (p *Provider) OpenDatabaseConnection() error {
	path := p.ProviderConfig.SqlitePath
	if path == "" {
		path = "tmp/example.db"
	}

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return fmt.Errorf("error opening sqlite connection %s: %w", path, err)
	}
	p.db = db
	return nil
}

func (p *Provider) CloseDatabaseConnection() error {
	err := p.db.Close()
	if err != nil {
		return fmt.Errorf("error closing sqlite connection %w", err)
	}
	return nil
}

func (p *Provider) Setup() error {
	err := p.OpenDatabaseConnection()
	if err != nil {
		return err
	}

	_, err = p.db.Exec(`CREATE TABLE IF NOT EXISTS secrets (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		key VARCHAR(64) UNIQUE NOT NULL,
		value VARCHAR(64) NOT NULL
	)`)
	if err != nil {
		return fmt.Errorf("error creating secrets table: %w", err)
	}
	return nil
}

func NewProvider(c Config) (*Provider, error) {
	p := &Provider{
		ProviderConfig: c,
	}
	err := p.Setup()
	if err != nil {
		return nil, err
	}
	return p, nil
}
