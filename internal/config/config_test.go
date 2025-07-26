package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadYaml(t *testing.T) {
	cfg, err := LoadYaml("test.yaml")
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.EqualValues(t, cfg.Token, "abcdefg")
	assert.EqualValues(t, cfg.ListenedChannel, "abc")
}
