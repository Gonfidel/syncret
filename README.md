# Secret Provider Interface

This repository defines a `Provider` interface for managing secrets. It is designed to be implemented by various secret management backends such as AWS Secrets Manager, or SQLite.

```go
package main

import (
	"fmt"
	"log"
	"github.com/gonfidel/syncret/providers/local"
)

func main() {
	// Setup the provider
	provider, err := local.NewProvider(local.Config{})
	if err != nil {
		log.Fatalf("Error setting up provider: %v", err)
	}

	// Set a secret
	err = provider.Set("exampleKey", "exampleValue")
	if err != nil {
		log.Fatalf("Error setting secret: %v", err)
	}

	// Get the secret
	secret, err := provider.Get("exampleKey")
	if err != nil {
		log.Fatalf("Error getting secret: %v", err)
	}
	fmt.Printf("Secret: %s\n", secret)

	// Check if the secret exists
	exists, err := provider.Exists("exampleKey")
	if err != nil {
		log.Fatalf("Error checking if secret exists: %v", err)
	}
	fmt.Printf("Secret exists: %v\n", exists)

	// Destroy the secret
	err = provider.Destroy("exampleKey")
	if err != nil {
		log.Fatalf("Error destroying secret: %v", err)
	}

	// Check if the secret still exists after destruction
	exists, err = provider.Exists("exampleKey")
	if err != nil {
		log.Fatalf("Error checking if secret exists: %v", err)
	}
	fmt.Printf("Secret exists after destruction: %v\n", exists)
}
```
