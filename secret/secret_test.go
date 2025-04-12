package secret_test

import (
	"testing"

	"github.com/gonfidel/syncret/internal/providers/mock"
	"github.com/gonfidel/syncret/secret"
)

func TestSecretManager_SetGetDestroy(t *testing.T) {
	manager, err := secret.NewMockManager(mock.Config{})
	if err != nil {
		t.Fatalf("failed to generate new provider: %v", err)
	}

	key := "foo/bar"
	value := "s3cr3t"

	if err = manager.Set(key, value); err != nil {
		t.Fatalf("failed to set secret: %v", err)
	}

	got, err := manager.Get(key)
	if err != nil {
		t.Fatalf("failed to get secret: %v", err)
	}
	if got != value {
		t.Errorf("expected value %q, got %q", value, got)
	}

	exists, err := manager.Exists(key)
	if err != nil {
		t.Fatalf("failed to check existence: %v", err)
	}
	if !exists {
		t.Error("expected secret to exist")
	}

	if err = manager.Destroy(key); err != nil {
		t.Fatalf("failed to destroy secret: %v", err)
	}

	exists, _ = manager.Exists(key)
	if exists {
		t.Error("expected secret to be deleted")
	}
}
