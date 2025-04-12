package local_test

import (
	"path/filepath"
	"testing"

	"github.com/gonfidel/syncret/providers/local"
)

func TestLocalProvider_SetGetDestroy(t *testing.T) {
	path := "test.db"
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, path)

	provider, err := local.NewProvider(local.Config{SqlitePath: dbPath})
	if err != nil {
		t.Fatalf("failed to create local provider: %v", err)
	}

	key := "TestKey"
	value := "HelloWorld!"

	if err := provider.Set(key, value); err != nil {
		t.Fatalf("failed to set secret: %v", err)
	}

	got, err := provider.Get(key)
	if err != nil {
		t.Fatalf("failed to get secret: %v", err)
	}
	if got != value {
		t.Errorf("expected %q, got %q", value, got)
	}

	if err := provider.Destroy(key); err != nil {
		t.Fatalf("failed to destroy secret: %v", err)
	}

	exists, err := provider.Exists(key)
	if err != nil {
		t.Fatalf("failed to check existence: %v", err)
	}
	if exists {
		t.Errorf("expected secret to not exist after destroy")
	}

	_, err = provider.Get(key)
	if err == nil {
		t.Errorf("expected error when getting deleted secret, got nil")
	}

	t.Cleanup(func() {
		err := provider.CloseDatabaseConnection()
		if err != nil {
			t.Errorf("failed closing local provider sqlite connection")
		}
	})
}
