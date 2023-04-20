package config

import (
	"log"

	"github.com/hashicorp/hcl/v2/hclsimple"
)

type Configuration struct {
}

func Load() (*Configuration, error) {
	var config Configuration
	err := hclsimple.DecodeFile("config.hcl", nil, &config)
	if err != nil {
		log.Fatalf("Failed to load configuration: %s", err)
	}
	return &config, nil
}
