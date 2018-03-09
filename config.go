package main

import (
	edgegrid "github.com/RafPe/go-edgegrid"
)

type Config struct {
	edgerc  string
	section string
}

func (c *Config) Client() (*edgegrid.Client, error) {

	clientOpts := edgegrid.ClientOptions{
		ConfigPath:    c.edgerc,
		ConfigSection: c.section,
	}

	client := edgegrid.NewClient(nil, &clientOpts)

	return client, nil
}
