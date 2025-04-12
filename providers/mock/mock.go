package mock

import (
	"fmt"
	"sync"
)

type Config struct{}

type Provider struct {
	secrets map[string]string
	mu      sync.RWMutex
	config  Config
}

func NewProvider(c Config) (*Provider, error) {
	p := Provider{
		secrets: make(map[string]string),
		config:  c,
	}
	err := p.Setup()
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (m *Provider) Setup() error {
	return nil
}

func (m *Provider) Set(key, value string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.secrets[key] = value
	return nil
}

func (m *Provider) Get(key string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	value, ok := m.secrets[key]
	if !ok {
		return "", fmt.Errorf("secret not found: %s", key)
	}
	return value, nil
}

func (m *Provider) Destroy(key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.secrets[key]; !ok {
		return fmt.Errorf("secret not found: %s", key)
	}
	delete(m.secrets, key)
	return nil
}

func (m *Provider) Exists(key string) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, ok := m.secrets[key]
	return ok, nil
}
