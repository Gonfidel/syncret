package secret

import (
	p "github.com/gonfidel/syncret/providers"
)

// // TODO (ngeorge): Review if we can remove this struct export
// type Secret struct {
// 	Key string
// 	Value string
// }

type SecretManager struct {
	provider p.Provider
}

func NewSecretManager(provider p.Provider) *SecretManager {
	return &SecretManager{provider}
}

func (sm *SecretManager) Get(key string) (string, error) {
	return sm.provider.Get(key)
}

func (sm *SecretManager) Set(key, value string) error {
	return sm.provider.Set(key, value)
}

func (sm *SecretManager) Destroy(key string) error {
	return sm.provider.Destroy(key)
}

func (sm *SecretManager) Exists(key string) (bool, error) {
	return sm.provider.Exists(key)
}
