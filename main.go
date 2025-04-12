package main

import (
	"fmt"

	aws "github.com/gonfidel/syncret/providers/aws"
	local "github.com/gonfidel/syncret/providers/local"
	// secret "github.com/gonfidel/syncret/secret"
)

func main() {
	// TODO (ngeorge): Review if I can rmove the secrets manager object in favor of the
	// provider. Feels redundant considering we can call the provider directly
	awsProvider, err := aws.NewProvider(aws.Config{})
	if err != nil {
		fmt.Printf("unable to build provider for aws secrets manager")
	}

	// awsVault := secret.NewSecretManager(awsProvider)

	err = awsProvider.Set("sample/test/123", "test")
	if err != nil {
		fmt.Printf("Err: %v", err)
	}
	awsTest, err := awsProvider.Get("sample/test/123")
	if err != nil {
		fmt.Printf("Err: %v", err)
	}
	fmt.Println(awsTest)

	err = awsProvider.Destroy("sample/test/123")
	if err != nil {
		fmt.Printf("failed to destroy secret: %v\n", err)
	}

	exists, err := awsProvider.Exists("sample/test/123")
	if err != nil {
		fmt.Printf("Err: %v\n", err)
	}
	fmt.Println(exists)

	awsTest, err = awsProvider.Get("sample/test/123")
	if err != nil {
		fmt.Printf("Err: %v", err)
	}
	fmt.Println(awsTest)

	////////////////////////////////
	localProvider, err := local.NewProvider(local.Config{})
	if err != nil {
		fmt.Printf("Err: %v", err)
	}
	// localVault := secret.NewSecretManager(localProvider)
	err = localProvider.Set("test2", "hello_world2")
	if err != nil {
		fmt.Printf("Err: %v\n", err)
	}

	v, err := localProvider.Get("test2")
	if err != nil {
		fmt.Printf("Err: %v\n", err)
	}
	fmt.Println(v)

	err = localProvider.Destroy("test2")
	if err != nil {
		fmt.Printf("Err: %v\n", err)
	}

	exists, err = localProvider.Exists("test2")
	if err != nil {
		fmt.Printf("Err: %v\n", err)
	}
	fmt.Println(exists)

	v, err = localProvider.Get("test2")
	if err != nil {
		fmt.Printf("Err: %v\n", err)
	}
	fmt.Println(v)
}
