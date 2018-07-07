package config

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfiguration_NotFoundError(t *testing.T) {
	_, err := LoadConfiguration()
	assert.Error(t, err)
}
