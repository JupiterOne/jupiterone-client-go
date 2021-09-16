package example

import (
	"fmt"

	client "github.com/jupiterone/jupiterone-client-go/jupiterone"
)

func main() {

	// Set configuration
	config := jupiterone.client.JupiterOneClientConfig{
		APIKey:    "api_key",
		AccountID: "accountID",
		Region:    "dev",
	}

	//Initialize client
	client, err := config.Client()

	if err != nil {
		fmt.Println("failed to create JupiterOne client: %s", err.Error())
	}

	//Do stuffs
	fmt.Print(client)
}
