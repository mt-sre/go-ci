package gocmd

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestModule(t *testing.T) {
	t.Parallel()

	gocmd, err := NewGoCmd()
	require.NoError(t, err)

	module, err := gocmd.Module(context.Background())
	require.NoError(t, err)

	assert.Equal(t, "github.com/mt-sre/go-ci", module)
}
