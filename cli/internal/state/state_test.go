package state_test

import (
	"testing"

	"github.com/lkmliz/mc-server/cli/internal/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad_MissingFile(t *testing.T) {
	s, err := state.Load(t.TempDir())
	require.NoError(t, err)
	assert.Nil(t, s)
}

func TestSave_CanBeLoaded(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, state.Save(dir, savedState()))
	loaded, err := state.Load(dir)
	require.NoError(t, err)
	require.NotNil(t, loaded)
	assert.Equal(t, savedState(), loaded)
}

func TestClear_MakesLoadReturnNil(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, state.Save(dir, savedState()))
	require.NoError(t, state.Clear(dir))
	s, err := state.Load(dir)
	require.NoError(t, err)
	assert.Nil(t, s)
}

func TestClear_MissingFileIsNotAnError(t *testing.T) {
	err := state.Clear(t.TempDir())
	assert.NoError(t, err)
}

func savedState() *state.State {
	return &state.State{Profile: "skyblock", InstanceIP: "1.2.3.4"}
}
