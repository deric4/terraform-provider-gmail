package gmail

import (
	"context"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var (
	contextTimeout = 15 * time.Second
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"gmail_label": resourceLabel(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	var c Config
	if err := c.loadAndValidate(); err != nil {
		return nil, err
	}

	return &c, nil
}

func contextWithTimeout() (context.Context, func()) {
	return context.WithTimeout(context.Background(), contextTimeout)
}
