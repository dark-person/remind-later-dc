// Package config is to load config.yaml to golang struct for discord bot usage.
package config

import (
	"fmt"
	"os"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

// Config for discord bots.
type DiscordConfig struct {
	// Token of discord bot. This token is required to send message to discord channel.
	Token string `koanf:"discord-token"`

	// Channel ID to listen
	ListenedChannel string `koanf:"discord-channel"`
}

// Load yaml config from given path,
// while no koanf instance will preserved (i.e. every call overwrite previous call).
//
// If failed to load config, then nil config will be returned with error.
func LoadYaml(path string) (*DiscordConfig, error) {
	var k = koanf.New(".")

	// Check if file exist
	if _, err := os.Stat(path); err != nil {
		return nil, fmt.Errorf("path %s does not exist", path)
	}

	// Start Load file
	err := k.Load(file.Provider(path), yaml.Parser())
	if err != nil {
		return nil, err
	}

	// Unmarshal to struct
	var out DiscordConfig
	err = k.Unmarshal("", &out)
	if err != nil {
		return nil, err
	}

	return &out, nil
}
