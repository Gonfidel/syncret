package local_test

import (
	"path/filepath"
	"testing"

	"github.com/gonfidel/syncret/internal/providers/local"
)

const (
	validEncryptionKey   = "12345678901234567890123456789012" // 32 bytes
	invalidEncryptionKey = "1234"                             // invalid
)

func setupProvider(t *testing.T, key string) *local.Provider {
	t.Helper()

	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")
	provider, err := local.NewProvider(local.Config{
		SqlitePath:    dbPath,
		EncryptionKey: key,
	})
	if err != nil {
		t.Fatalf("failed to create local provider: %v", err)
	}

	t.Cleanup(func() {
		if err = provider.CloseDatabaseConnection(); err != nil {
			t.Errorf("failed closing local provider sqlite connection")
		}
	})

	return provider
}

func TestLocalProvider_SetGetDestroy(t *testing.T) {
	provider := setupProvider(t, validEncryptionKey)

	const key = "TestKey"
	const value = "HelloWorld!"

	t.Run("Set", func(t *testing.T) {
		if err := provider.Set(key, value); err != nil {
			t.Fatalf("failed to set secret: %v", err)
		}
	})

	t.Run("Get", func(t *testing.T) {
		got, err := provider.Get(key)
		if err != nil {
			t.Fatalf("failed to get secret: %v", err)
		}
		if got != value {
			t.Errorf("expected %q, got %q", value, got)
		}
	})

	t.Run("Exists", func(t *testing.T) {
		exists, err := provider.Exists(key)
		if err != nil {
			t.Fatalf("failed to check existence: %v", err)
		}
		if !exists {
			t.Errorf("expected secret to exist")
		}
	})

	t.Run("Destroy", func(t *testing.T) {
		if err := provider.Destroy(key); err != nil {
			t.Fatalf("failed to destroy secret: %v", err)
		}

		exists, err := provider.Exists(key)
		if err != nil {
			t.Fatalf("failed to check existence after destroy: %v", err)
		}
		if exists {
			t.Errorf("expected secret to not exist after destroy")
		}

		_, err = provider.Get(key)
		if err == nil {
			t.Errorf("expected error when getting deleted secret, got nil")
		}
	})
}

func TestLocalProvider_InvalidEncryptionKey(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	_, err := local.NewProvider(local.Config{
		SqlitePath:    dbPath,
		EncryptionKey: invalidEncryptionKey,
	})

	if err == nil {
		t.Errorf("expected error when creating provider with invalid encryption key")
	}
}
