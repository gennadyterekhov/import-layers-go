package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCanGetConfigFromFile(t *testing.T) {
	cfg, err := FromFile()
	require.NoError(t, err)

	assert.Equal(t, 6, cfg.GetLayer("github.com/gennadyterekhov/import-layers-go/internal/project"))
	assert.Equal(t, 5, cfg.GetLayer("github.com/gennadyterekhov/import-layers-go/internal/data"))
	assert.Equal(t, 4, cfg.GetLayer("github.com/gennadyterekhov/import-layers-go/internal/config"))
	assert.Equal(t, 3, cfg.GetLayer("github.com/gennadyterekhov/import-layers-go/internal/reporter"))
	assert.Equal(t, 2, cfg.GetLayer("github.com/gennadyterekhov/import-layers-go/internal/filechecker"))
	assert.Equal(t, 1, cfg.GetLayer("github.com/gennadyterekhov/import-layers-go/internal/analyzer"))
}
