package main

import (
	"github.com/deric4/terraform-provider-gmail/gmail"
	"github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/terraform"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return gmail.Provider()
		},
	})
}
