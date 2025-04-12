package mock_test

import (
	"testing"

	"github.com/gonfidel/syncret/internal/providers/mock"
)

func setupProvider(t *testing.T) *mock.Provider {
	t.Helper()

	provider, err := mock.NewProvider(mock.Config{})
	if err != nil {
		t.Fatalf("failed to create mock provider: %v", err)
	}

	return provider
}

func TestMockProvider_SetGetDestroy(t *testing.T) {
	provider := setupProvider(t)

	const key = "foo/bar"
	const value = "s3cr3t"

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
			t.Errorf("expected value %q, got %q", value, got)
		}
	})

	t.Run("Exists", func(t *testing.T) {
		exists, err := provider.Exists(key)
		if err != nil {
			t.Fatalf("failed to check existence: %v", err)
		}
		if !exists {
			t.Error("expected secret to exist")
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
			t.Error("expected secret to be deleted")
		}

		_, err = provider.Get(key)
		if err == nil {
			t.Error("expected error when getting deleted secret, got nil")
		}
	})
}
