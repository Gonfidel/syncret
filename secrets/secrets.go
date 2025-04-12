package secrets

type Store interface {
	Destroy(key string) error
	Exists(key string) (bool, error)
	Get(key string) (string, error)
	Init() error
	Set(key, value string) error
	Shutdown() error
}
