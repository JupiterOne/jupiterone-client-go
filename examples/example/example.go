package main

import (
	"fmt"

	client "github.com/jupiterone/jupiterone-client-go/jupiterone"
)

func main() {
	var entityProps client.EntityProperties

	// Set configuration
	config := client.JupiterOneClientConfig{
		APIKey:    "key",
		AccountID: "j1dev",
		Region:    "dev",
	}

	entityProps.Key = "go-client-key"
	entityProps.Type = "go_client_type"
	entityProps.Class = "Record"

	//Initialize client
	client, err := config.Client()

	if err != nil {
		fmt.Println("failed to create JupiterOne client: %s", err.Error())
	}

	//Do stuffs
	// fmt.Print(client)
	fmt.Print(client.CreateEntity(entityProps))
}
