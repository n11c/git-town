package test

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGitManagerCreateMemoizedEnvironment(t *testing.T) {
	dir := createTempDir(t)
	gm := NewGitManager(dir)

	err := gm.CreateMemoizedEnvironment()

	assert.Nil(t, err, "creating memoized environment failed")
	memoizedPath := path.Join(dir, "memoized")
	_, err = os.Stat(memoizedPath)
	assert.Falsef(t, os.IsNotExist(err), "memoized directory %q not found", memoizedPath)
}

func TestGitManagerCreateScenarioEnvironment(t *testing.T) {
	dir := createTempDir(t)
	gm := NewGitManager(dir)
	err := gm.CreateMemoizedEnvironment()
	assert.Nil(t, err, "creating memoized environment failed")

	result, err := gm.CreateScenarioEnvironment("foo")

	assert.Nil(t, err, "cannot create scenario environment")
	_, err = os.Stat(result.DeveloperRepo.Dir)
	assert.False(t, os.IsNotExist(err), "scenario environment directory %q not found", result.DeveloperRepo.Dir)
}
