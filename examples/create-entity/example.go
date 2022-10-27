package main

import (
	"fmt"
	"os"

	j1 "github.com/jupiterone/jupiterone-client-go/jupiterone"
)

func getEnvWithDefault(key string, defaultVal string) string {
    value, exists := os.LookupEnv(key)
    if !exists {
        value = defaultVal
    }
    return value
}

func main() {
	// Set configuration
	config := j1.Config{
		APIKey:    getEnvWithDefault("J1_API_TOKEN", ""),
		AccountID: getEnvWithDefault("J1_ACCOUNT", ""),
		Region:    getEnvWithDefault("J1_REGION", "us"),
	}

	entityProps := j1.EntityProperties{
		Key:   "go-client-key",
		Type:  "go_client_type",
		Class: "Record",
		Properties: map[string]interface{}{
			"displayName": "exampleRecord",
			"stringVal": "Mississippi",
			"client": "jupiterone-client-go",
			"isBool": true,
	    },
	}

	// Initialize client
	client, err := j1.NewClient(&config)
	if err != nil {
		fmt.Printf("failed to create JupiterOne client: %s", err.Error())
	}

	entity, err := client.Entity.Create(entityProps)
	if err != nil {
		fmt.Printf("failed to create entity: %s", err.Error())
	}

	fmt.Print(entity)
}
