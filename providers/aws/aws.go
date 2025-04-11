package aws

import (
	"errors"
)

type Config struct {
	AwsSecretKeyId string
	AwsSecretAccessKey string
}

type AwsProvider struct {
	Config Config
}

func (p *AwsProvider) Get(key string) (string, error) {
	return "", errors.New("not implemented")
}
func (p *AwsProvider) Set(key, value string) error {
	return errors.New("not implemented")
}
func (p *AwsProvider) Destroy(key string) error {
	return errors.New("not implemented")
}

func (p *AwsProvider) Setup() {}

func NewProvider(c Config) *AwsProvider {
	return &AwsProvider{
		Config: c,
	}
}
