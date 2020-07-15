package main

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultCache(t *testing.T) {
	_, err := DefaultCache()
	assert.NoError(t, err)
}

func TestIsHitMiss(t *testing.T) {
	cache, err := DefaultCache()
	assert.NoError(t, err)
	hit := cache.IsHit("github.com/this-user-doesn't-exist/not-a-repository-that-exists")
	assert.False(t, hit)
}

func TestCacheConfig(t *testing.T) {
	_, filename, _, ok := runtime.Caller(0)
	assert.Truef(t, ok, "failed to get filename of test file cache_test.go")
	parent := filepath.Dir(filename)
	testDataSlash := fmt.Sprintf("%s/testdata", parent)
	cache := &Cache{filepath.FromSlash(testDataSlash)}
	conf, err := cache.Config()
	assert.NoError(t, err)
	assert.Equal(t, conf.License, "GPLv3")
	assert.Equal(t, conf.Name, "YourName/dotfiles")
}

func TestCleanDoesntExist(t *testing.T) {
	invalidFile := "this-isn't-a-present-file-please-don't-let-this-be-a-present-file.abcdefgh"
	cache := &Cache{invalidFile}
	err := cache.Clean()
	assert.NoError(t, err)
}
