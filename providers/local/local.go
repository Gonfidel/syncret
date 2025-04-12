package local

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"runtime"

	_ "modernc.org/sqlite" // Driver necessary but never directly called

	"github.com/gonfidel/syncret/secrets"
)

type Config struct {
	SqlitePath    string
	EncryptionKey string
}

type Provider struct {
	ProviderConfig      Config
	db                  *sql.DB
	encryptionByteArray []byte
}

func NewProvider(c Config) (secrets.Store, error) {
	p := &Provider{ProviderConfig: c}
	if err := p.Init(); err != nil {
		return nil, err
	}

	runtime.SetFinalizer(p, func(p *Provider) {
		_ = p.Shutdown()
	})

	return p, nil
}

func (p *Provider) Init() error {
	err := p.validateEncryptionKey()
	if err != nil {
		return err
	}

	err = p.OpenDatabaseConnection()
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

func (p *Provider) Shutdown() error {
	err := p.CloseDatabaseConnection()
	return err
}

func (p *Provider) Get(key string) (string, error) {
	var encryptedValue string
	row := p.db.QueryRow(`SELECT value FROM secrets WHERE key = ?;`, key)
	err := row.Scan(&encryptedValue)

	if err != nil {
		return "", fmt.Errorf("unable to get secret with key \"%s\" %w", key, err)
	}

	decryptedValue, err := p.decrypt(encryptedValue)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt secret with key \"%s\": %w", key, err)
	}

	return decryptedValue, nil
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
	encryptedValue, err := p.encrypt(value)
	if err != nil {
		return err
	}

	_, err = p.db.Exec(`
		INSERT INTO secrets (key, value)
		VALUES (?, ?)
		ON CONFLICT(key) DO UPDATE SET value = excluded.value;
	`, key, encryptedValue)
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

func (p *Provider) encrypt(plaintext string) (string, error) {
	block, err := aes.NewCipher(p.encryptionByteArray)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))

	iv := make([]byte, aes.BlockSize)
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return "", fmt.Errorf("failed to generate IV: %w", err)
	}

	copy(ciphertext[:aes.BlockSize], iv)

	stream := cipher.NewCFBEncrypter(block, iv) // #nosec G407 -- IV is randomized above
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(plaintext))

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (p *Provider) decrypt(ciphertext string) (string, error) {
	block, err := aes.NewCipher(p.encryptionByteArray)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	ciphertextBytes, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", fmt.Errorf("failed to decode ciphertext: %w", err)
	}

	if len(ciphertextBytes) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}
	iv := ciphertextBytes[:aes.BlockSize]
	ciphertextBytes = ciphertextBytes[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertextBytes, ciphertextBytes)

	return string(ciphertextBytes), nil
}

func (p *Provider) validateEncryptionKey() error {
	keyLen := len(p.ProviderConfig.EncryptionKey)
	if keyLen != 16 && keyLen != 24 && keyLen != 32 {
		return fmt.Errorf("invalid encryption key length: got %d bytes, want 16, 24, or 32", keyLen)
	}
	p.encryptionByteArray = []byte(p.ProviderConfig.EncryptionKey)
	return nil
}
