package main

import (
	"log"
	"os"

	"github.com/jupiterone/jupiterone-client-go/internal/graphql"
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

	client, err := j1.NewClient(&config)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}

	results, err := client.Query.Query(j1.QueryInput{
		Query: "FIND jupiterone_integration AS x RETURN x.id",
	})

	if err != nil {
		log.Fatalf("failed query output: %v", err)
	}

	for _, result := range results.([]map[string]interface{}) {
		integrationInstanceId := result["x.id"].(string)
		integrationInstance, err := client.Integration.GetIntegrationInstance(integrationInstanceId)
		if err != nil {
			log.Printf("failed query output: %v", err)
			continue
		}

		if tags, ok := integrationInstance.IntegrationInstance.Config["@tag"]; ok {
			tags.(map[string]interface{})["j1.sourcefilter"] = "layer0"
		}

		input := graphql.UpdateIntegrationInstanceInput{
			Name:                        integrationInstance.IntegrationInstance.Name,
			SourceIntegrationInstanceId: integrationInstance.IntegrationInstance.SourceIntegrationInstanceId,
			PollingInterval:             integrationInstance.IntegrationInstance.PollingInterval,
			PollingIntervalCronExpression: graphql.IntegrationPollingIntervalCronExpressionInput{
				DayOfWeek: integrationInstance.IntegrationInstance.PollingIntervalCronExpression.DayOfWeek,
				Hour:      integrationInstance.IntegrationInstance.PollingIntervalCronExpression.Hour,
			},
			Description:     integrationInstance.IntegrationInstance.Name,
			Config:          integrationInstance.IntegrationInstance.Config,
			OffsiteComplete: integrationInstance.IntegrationInstance.OffsiteComplete,
		}

		resp, err := client.Integration.UpdateIntegrationInstance(integrationInstanceId, input)
		if err != nil {
			log.Printf("failed to update integration instance: %v", err)
			continue
		}

		log.Printf("updated integration instance: %v\n", resp)
	}

	log.Print("program complete.")
}
