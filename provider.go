package main

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"edgerc": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"section": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
		ConfigureFunc: providerConfigure,
		ResourcesMap: map[string]*schema.Resource{
			"xakamai_network_list": resourceNetlist(),
		},
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		edgerc:  d.Get("edgerc").(string),
		section: d.Get("section").(string),
	}

	return config.Client()
}
