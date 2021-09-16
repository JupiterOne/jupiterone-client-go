package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/jupiterone/jupiterone-client-go/jupiterone"
)

//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs
func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: jupiterone.Provider})
}
