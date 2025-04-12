package mock

import (
	"fmt"
	"runtime"
	"sync"
)

type Config struct{}

type Provider struct {
	secrets        map[string]string
	mu             sync.RWMutex
	ProviderConfig Config
}

func NewProvider(c Config) (*Provider, error) {
	p := &Provider{
		ProviderConfig: c,
		secrets:        make(map[string]string),
	}

	if err := p.Init(); err != nil {
		return nil, err
	}

	runtime.SetFinalizer(p, func(p *Provider) {
		_ = p.Shutdown()
	})

	return p, nil
}

func (m *Provider) Init() error {
	return nil
}

func (m *Provider) Shutdown() error {
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
