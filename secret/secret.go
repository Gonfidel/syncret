package secret

import (
	p "github.com/gonfidel/syncret/providers"
)

type Manager struct {
	provider p.Provider
}

func NewManager(provider p.Provider) *Manager {
	return &Manager{provider}
}

func (sm *Manager) Get(key string) (string, error) {
	return sm.provider.Get(key)
}

func (sm *Manager) Set(key, value string) error {
	return sm.provider.Set(key, value)
}

func (sm *Manager) Destroy(key string) error {
	return sm.provider.Destroy(key)
}

func (sm *Manager) Exists(key string) (bool, error) {
	return sm.provider.Exists(key)
}
