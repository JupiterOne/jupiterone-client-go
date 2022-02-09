package main

import (
	"fmt"

	j1 "github.com/jupiterone/jupiterone-client-go/jupiterone"
)

func main() {
	// var entityProps j1.EntityProperties

	// Set configuration
	config := j1.Config{
		APIKey:    "api_key",
		AccountID: "accountId",
		Region:    "dev",
	}

	// entityProps.Key = "go-client-key"
	// entityProps.Type = "go_client_type"
	// entityProps.Class = "Record"

	entityProps := j1.EntityProperties{
		Key:   "go-client-key",
		Type:  "go_client_type",
		Class: "Record",
	}

	// Initialize client
	client, err := j1.NewClient(&config)
	if err != nil {
		fmt.Printf("failed to create JupiterOne client: %s", err.Error())
	}

	// Do stuffs
	// fmt.Print(client)
	entity, err := client.Entity.Create(entityProps)

	fmt.Print(entity)
}
