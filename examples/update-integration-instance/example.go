package main

import (
	"encoding/json"
	"log"
	"os"

	j1 "github.com/jupiterone/jupiterone-client-go/jupiterone"
	"github.com/jupiterone/jupiterone-client-go/jupiterone/graphql"
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

	resultsAsMap, ok := results.(map[string]interface{})
	if !ok {
		log.Fatalf("failed to convert results to map[string]interface{}")
	}

	data, ok := resultsAsMap["data"].([]interface{})
	if !ok {
		log.Fatalf("failed to convert data to []interface{}")
	}

	for _, dataRow := range data {

		valuesHash, ok := dataRow.(map[string]interface{})
		if !ok {
			log.Printf("failed to convert dataRow to map[string]interface{}")
			continue
		}

		integrationInstanceId, ok := valuesHash["x.id"].(string)
		if !ok {
			log.Printf("failed to retrieve integration instance id")
			continue
		}

		/*

			The following code is the core of this script.
			It follows the standard pattern of: get current
			state, update the state, and then save the state.

		*/

		integrationInstance, err := client.Integration.GetIntegrationInstance(integrationInstanceId)
		if err != nil {
			log.Printf("failed query output: %v", err)
			continue
		}

		configUpdate := map[string]interface{}{}
		haveOneComponentOfAWSIntegration := false

		accountId, ok := integrationInstance.IntegrationInstance.Config["accountId"]
		if ok {
			configUpdate["accountId"] = accountId
			haveOneComponentOfAWSIntegration = true
		}

		roleArn, ok := integrationInstance.IntegrationInstance.Config["roleArn"]
		if ok {
			configUpdate["roleArn"] = roleArn
			haveOneComponentOfAWSIntegration = true
		}

		// this is a check to ensure that if we have one
		// component of the tags that make up an aws integration,
		// that we have both. if we don't have both, we can find
		// ourselves in a sticky situation.
		if haveOneComponentOfAWSIntegration {
			if configUpdate["accountId"] == nil {
				log.Printf("accountId is missing from the config, bailing on this integration instance")
				continue
			}

			if configUpdate["roleArn"] == nil {
				log.Printf("roleArn is missing from the config, bailing on this integration instance")
				continue
			}
		}

		tags, ok := integrationInstance.IntegrationInstance.Config["@tag"]
		if !ok {
			log.Printf("unable to reach the @tag property for %v", integrationInstanceId)
			continue
		}

		tagsAsMap, ok := tags.(map[string]interface{})
		if !ok {
			log.Printf("unable to convert tags to map[string]interface{}")
			continue
		}

		/*

			If for example, here, we want to delete a tag, instead
			of adding a tag, we would use the following code:

			delete(tagsAsMap, "j1.sourcefilter")

		*/

		tagsAsMap["j1.sourcefilter"] = "layer0"
		configUpdate["@tag"] = tagsAsMap

		input := graphql.UpdateIntegrationInstanceInput{
			Config: &configUpdate,
		}

		resp, err := client.Integration.UpdateIntegrationInstance(integrationInstanceId, input)
		if err != nil {
			log.Printf("failed to update integration instance: %v", err)
			continue
		}

		var jsonOutput map[string]interface{}

		b, err := json.Marshal(*resp)
		if err != nil {
			log.Fatalf("failed to marshal output: %v", err)
		}

		err = json.Unmarshal(b, &jsonOutput)
		if err != nil {
			log.Fatal("failed to unmarshal output")
		}

		jsonString, err := json.MarshalIndent(jsonOutput, "", "  ")
		if err != nil {
			log.Printf("failed to marshal output: %v\n", err)
		}

		log.Println("updated integration instance")
		log.Println(string(jsonString))

	}

	log.Print("program complete.")
}
