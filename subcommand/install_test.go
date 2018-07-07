package subcommand

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInstall(t *testing.T) {
	res := Install()

	assert.Exactly(t, 1, res)
}
