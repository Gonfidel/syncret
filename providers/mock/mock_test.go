package mock_test

import (
	"testing"

	"github.com/gonfidel/syncret/providers/mock"
)

func TestMockProvider_SetGetDestroy(t *testing.T) {
	mockProvider, err := mock.NewProvider(mock.Config{})
	if err != nil {
		t.Fatalf("failed to generate new provider: %v", err)
	}

	key := "foo/bar"
	value := "s3cr3t"

	if err := mockProvider.Set(key, value); err != nil {
		t.Fatalf("failed to set secret: %v", err)
	}

	got, err := mockProvider.Get(key)
	if err != nil {
		t.Fatalf("failed to get secret: %v", err)
	}
	if got != value {
		t.Errorf("expected value %q, got %q", value, got)
	}

	exists, err := mockProvider.Exists(key)
	if err != nil {
		t.Fatalf("failed to check existence: %v", err)
	}
	if !exists {
		t.Error("expected secret to exist")
	}

	if err := mockProvider.Destroy(key); err != nil {
		t.Fatalf("failed to destroy secret: %v", err)
	}

	exists, _ = mockProvider.Exists(key)
	if exists {
		t.Error("expected secret to be deleted")
	}
}
