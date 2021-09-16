# JupiterOne Client Go

## Requirements

- [Go](https://golang.org/doc/install)

## Usage

```
package example

import (
	"fmt"

	"../jupiterone/client"
)

func main() {

	// Set configuration
	config := client.JupiterOneClientConfig{
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

```
