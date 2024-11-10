package config

import (
	"fmt"

	"gopkg.in/yaml.v2"
)

func RpConfigFromBytes(b []byte) (*reverseproxyConfig, error) {

	cfg := &reverseproxyConfig{}

	err := yaml.UnmarshalStrict(b, cfg)
	if err != nil {
		return nil, fmt.Errorf("Error un-marshaling config: %w", err)
	}

	return cfg, nil
}