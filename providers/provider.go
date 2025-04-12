package provider

type Provider interface {
	Setup() error
	Get(key string) (string, error)
	Set(key, value string) error
	Destroy(key string) error
	Exists(key string) (bool, error)
}