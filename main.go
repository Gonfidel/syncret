package main

import (
	"fmt"
	// aws "github.com/gonfidel/syncret/providers/aws"
	local "github.com/gonfidel/syncret/providers/local"
	secret "github.com/gonfidel/syncret/secret"
)

func main(){
	// awsProvider := aws.NewProvider(
	// 	aws.Config {
	// 		AwsSecretKeyId: "key_id",
	// 		AwsSecretAccessKey: "secret_key",
	// 	},
	// )

	// awsVault := secret.NewSecretManager(awsProvider)
	// awsTest, err := awsVault.Get("test")
	// if err != nil {
	// 	fmt.Printf("Err: %v", err)
	// }
	// fmt.Println(awsTest)

	////////////////////////////////
	localProvider := local.NewProvider(local.Config{})
	localVault := secret.NewSecretManager(localProvider)
	err := localVault.Set("test2", "hello_world2")
	if err != nil {
		fmt.Printf("Err: %v", err)
	}

	_, err = localVault.Get("test2")
	if err != nil {
		fmt.Printf("Err: %v", err)
	}

	err = localVault.Destroy("test2")
	if err != nil {
		fmt.Printf("Err: %v", err)
	}

	_, err = localVault.Get("test2")
	if err != nil {
		fmt.Printf("Err: %v", err)
	}
}