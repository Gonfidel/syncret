package mock

import (
	"fmt"
	"sync"
)

type Config struct{}

type MockProvider struct {
	secrets map[string]string
	mu      sync.RWMutex
}

func NewProvider(c Config) (*MockProvider, error) {
	p := MockProvider{
		secrets: make(map[string]string),
	}
	err := p.Setup()
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (m *MockProvider) Setup() error {
	return nil
}

func (m *MockProvider) Set(key, value string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.secrets[key] = value
	return nil
}

func (m *MockProvider) Get(key string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	value, ok := m.secrets[key]
	if !ok {
		return "", fmt.Errorf("secret not found: %s", key)
	}
	return value, nil
}

func (m *MockProvider) Destroy(key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.secrets[key]; !ok {
		return fmt.Errorf("secret not found: %s", key)
	}
	delete(m.secrets, key)
	return nil
}

func (m *MockProvider) Exists(key string) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, ok := m.secrets[key]
	return ok, nil
}
