package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/NickHackman/dots/config"
	"golang.org/x/tools/go/vcs"
)

// Cache is a wrapper around the path to the cache directory.
//
// The default case the Dir will be `XDG_CACHE_HOME/dots` or `~/.cache/dots`
type Cache struct {
	Dir string
}

// DefaultCache creates a Cache instance using the default value of
// `XDG_CACHE_HOME/dots`
func DefaultCache() (*Cache, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return nil, fmt.Errorf("failed to find cache directory: %w", err)
	}
	dir := filepath.Join(cacheDir, "dots")
	if err = os.MkdirAll(dir, os.ModeDir); err != nil {
		return nil, fmt.Errorf("failed to mkdir %s: %w", cacheDir, err)
	}
	return &Cache{dir}, nil
}

// IsHit checks to see if the cache already has a repository downloaded locally
//
// path is expected to be of the form `$domain/$username/$repoName`
// for example: `github.com/NickHackman/dotfiles`
func (cache *Cache) IsHit(path string) bool {
	repoPath := filepath.Join(cache.Dir, path)
	file, err := os.Stat(repoPath)
	if os.IsNotExist(err) || err != nil {
		return false
	}
	if !file.IsDir() {
		return false
	}

	dotsPath, err := config.FindConfig(repoPath)
	if err != nil {
		return false
	}

	if _, err = os.Stat(dotsPath); os.IsNotExist(err) || err != nil {
		return false
	}
	return true
}

// Config gets the dots configuration file at a given repository path,
// repository path is expected to be of the form `$domain/$username/$repoName`.
//
// Config cannot on its own determine if the cache is present, prior to calling `Config`
// one should call `IsHit` to verify the existence of the cache.
func (cache *Cache) Config() (*config.DotsConfig, error) {
	dotsConfig, err := config.Parse(cache.Dir)
	if err != nil {
		return nil, err
	}
	return dotsConfig, nil
}

// Clean completely removes all sub directories of `Cache.Dir`
func (cache *Cache) Clean() error {
	return os.RemoveAll(cache.Dir)
}

// UpgradeAll upgrades all repositories it finds inside of Cache
// in a breadth first fashion.
func (cache *Cache) UpgradeAll() error {
	queue := []string{cache.Dir}
	for len(queue) != 0 {
		front := queue[0]

		info, err := os.Stat(front)
		if err != nil {
			return err
		}

		if info.IsDir() {
			files, err := ioutil.ReadDir(front)
			if err != nil {
				return err
			}

			for _, file := range files {
				path := filepath.Join(front, file.Name())
				queue = append(queue, path)
			}
		}

		if err = cache.Upgrade(front); err != nil {
			return err
		}
	}
	return nil
}

// Upgrade upgrades a specific repository expected to exist at path
func (cache *Cache) Upgrade(path string) error {
	vcs, root, err := vcs.FromDir(path, path)
	if err != nil {
		return err
	}

	return vcs.Download(root)
}
