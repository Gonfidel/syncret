package aws

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	smTypes "github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
)

type Config struct{}

type Provider struct {
	ProviderConfig Config
	awsConfig      config.Config
	smClient       *secretsmanager.Client
}

func (p *Provider) Get(key string) (string, error) {
	resp, err := p.smClient.GetSecretValue(context.Background(), &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(key),
	})
	if err != nil {
		return "", fmt.Errorf("failed to get secret value with key \"%s\": %w", key, err)
	}

	if resp.SecretString != nil {
		return *resp.SecretString, nil
	}
	if resp.SecretBinary != nil {
		decoded := base64.StdEncoding.EncodeToString(resp.SecretBinary)
		return decoded, nil
	}
	return "", errors.New("unknown secret format")
}

func (p *Provider) Exists(key string) (bool, error) {
	_, err := p.smClient.DescribeSecret(context.Background(), &secretsmanager.DescribeSecretInput{
		SecretId: aws.String(key),
	})
	if err != nil {
		var notFound *smTypes.ResourceNotFoundException
		if ok := errors.As(err, &notFound); ok {
			// Secret doesn't exist
			return false, nil
		}
		// Some other error occurred
		return false, fmt.Errorf("failed to describe secret: %w", err)
	}
	// Secret exists
	return true, nil
}

func (p *Provider) Set(key, value string) error {
	exists, err := p.Exists(key)
	if err != nil {
		return err
	}
	if exists {
		_, err = p.smClient.UpdateSecret(context.Background(), &secretsmanager.UpdateSecretInput{
			SecretId:     aws.String(key),
			SecretString: aws.String(value),
		})
		if err != nil {
			return fmt.Errorf("unable to set secrets %s: %w", key, err)
		}
	} else {
		_, err = p.smClient.CreateSecret(context.Background(), &secretsmanager.CreateSecretInput{
			Name:         aws.String(key),
			SecretString: aws.String(value),
		})
		if err != nil {
			return fmt.Errorf("unable to set secrets %s: %w", key, err)
		}
	}

	return nil
}

func (p *Provider) Destroy(key string) error {
	exists, err := p.Exists(key)
	if err != nil {
		return err
	}

	if exists {
		_, err = p.smClient.DeleteSecret(context.Background(), &secretsmanager.DeleteSecretInput{
			ForceDeleteWithoutRecovery: aws.Bool(true),
			SecretId:                   aws.String(key),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Provider) Setup() error {
	awsConfig, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return err
	}
	awsSecretsManager := secretsmanager.NewFromConfig(awsConfig)
	p.awsConfig = awsConfig
	p.smClient = awsSecretsManager
	return nil
}

func NewProvider(c Config) (*Provider, error) {
	p := &Provider{
		ProviderConfig: c,
	}
	err := p.Setup()
	if err != nil {
		return nil, err
	}
	return p, nil
}
