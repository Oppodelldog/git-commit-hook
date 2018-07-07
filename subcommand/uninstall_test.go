package subcommand

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUninstall(t *testing.T) {
	res := Uninstall()

	assert.Exactly(t, 1, res)
}
