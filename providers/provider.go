package provider

// import (
// 	s "github.com/gonfidel/syncret/secret"
// )

type Provider interface {
	Setup()
	Get(key string) (string, error)
	Set(key, value string) error
	Destroy(key string) error
}