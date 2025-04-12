package secret

import (
	aws "github.com/gonfidel/syncret/internal/providers/aws"
	local "github.com/gonfidel/syncret/internal/providers/local"
	mock "github.com/gonfidel/syncret/internal/providers/mock"
)

type Store interface {
	Init() error
	Shutdown() error
	Get(key string) (string, error)
	Set(key, value string) error
	Destroy(key string) error
	Exists(key string) (bool, error)
}

type Manager struct {
	secretStore Store
}

func NewLocalManager(c local.Config) (*Manager, error) {
	secretStore, err := local.NewProvider(c)
	if err != nil {
		return nil, err
	}

	return &Manager{secretStore}, nil
}

func NewAwsManager(c aws.Config) (*Manager, error) {
	secretStore, err := aws.NewProvider(c)
	if err != nil {
		return nil, err
	}

	return &Manager{secretStore}, nil
}

func NewMockManager(c mock.Config) (*Manager, error) {
	secretStore, err := mock.NewProvider(c)
	if err != nil {
		return nil, err
	}

	return &Manager{secretStore}, nil
}

func (sm *Manager) Shutdown() error {
	return sm.secretStore.Shutdown()
}

func (sm *Manager) Get(key string) (string, error) {
	return sm.secretStore.Get(key)
}

func (sm *Manager) Set(key, value string) error {
	return sm.secretStore.Set(key, value)
}

func (sm *Manager) Destroy(key string) error {
	return sm.secretStore.Destroy(key)
}

func (sm *Manager) Exists(key string) (bool, error) {
	return sm.secretStore.Exists(key)
}
