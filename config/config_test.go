package config_test

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/NickHackman/dots/config"
	"github.com/stretchr/testify/assert"
)

// Get the path to the testdata directory for testing
func pathToTestData() (string, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", errors.New("failed to get filename of test file config_test.go")
	}
	parent := filepath.Dir(filename)
	root := filepath.Dir(parent)
	testDataSlash := fmt.Sprintf("%s/testdata", root)
	return filepath.FromSlash(testDataSlash), nil
}

func TestValidParse(t *testing.T) {
	testData, err := pathToTestData()
	assert.NoErrorf(t, err, "failed to setup config_test.go testing: %w", err)

	configDir, err := os.UserConfigDir()
	assert.NoErrorf(t, err, "failed to setup config_test.go testing can't locate `$XDG_CONFIG_HOME`: %w", err)

	homeDir, err := os.UserHomeDir()
	assert.NoErrorf(t, err, "failed to setup config_test.go testing can't locate `$HOME`: %w", err)

	tests := []struct {
		path     string
		expected *config.DotsConfig
	}{
		{
			path: "template.yml",
			expected: &config.DotsConfig{
				Name:    "YourName/dotfiles",
				License: "GPLv3",
				URL:     "https://github.com/NickHackman/dots",
				Dotfiles: []config.Dotfile{
					{
						Name:        "bspwm",
						Description: "A simple configuration file for the Binary Space Partition Window Manager",
						Source:      fmt.Sprintf("%s%cbspwm", testData, os.PathSeparator),
						Destination: fmt.Sprintf("%s%cbspwm", configDir, os.PathSeparator),
					},
					{
						Name:            "keybinds",
						Description:     "Keybindings that escape <-> capslock and handle function keys",
						Source:          fmt.Sprintf("%s%ckeybinds", testData, os.PathSeparator),
						Destination:     homeDir,
						InstallChildren: true,
					},
				},
			},
		},
		{
			path: "no-dotfiles.yml",
			expected: &config.DotsConfig{
				Name:    "YourName/dotfiles",
				License: "GPLv3",
				URL:     "https://github.com/NickHackman/dots",
			},
		},
		{
			path: "empty-dotfiles.yml",
			expected: &config.DotsConfig{
				Name:    "YourName/dotfiles",
				License: "GPLv3",
				URL:     "https://github.com/NickHackman/dots",
			},
		},
		{
			path: "many-dotfiles.yml",
			expected: &config.DotsConfig{
				Name:    "YourName/dotfiles",
				License: "GPLv3",
				URL:     "https://github.com/NickHackman/dots",
				Dotfiles: []config.Dotfile{
					{
						Name:        "bspwm",
						Description: "A simple configuration file for the Binary Space Partition Window Manager",
						Source:      fmt.Sprintf("%s%cbspwm", testData, os.PathSeparator),
						Destination: fmt.Sprintf("%s%cbspwm", configDir, os.PathSeparator),
					},
					{
						Name:            "keybinds",
						Description:     "Keybindings that escape <-> capslock and handle function keys",
						Source:          fmt.Sprintf("%s%ckeybinds", testData, os.PathSeparator),
						Destination:     homeDir,
						InstallChildren: true,
					},
					{
						Name:        "test1",
						Description: "description",
						Source:      fmt.Sprintf("%s%ctest1", testData, os.PathSeparator),
						Destination: fmt.Sprintf("%s%ctest1", configDir, os.PathSeparator),
					},
				},
			},
		},
	}

	for _, test := range tests {
		fullPathSlash := filepath.Join(testData, test.path)
		fullPath := filepath.FromSlash(fullPathSlash)

		t.Run(test.path, func(t *testing.T) {
			DotsConfig, err := config.ParseFile(fullPath)
			assert.NoError(t, err)
			assert.Equal(t, test.expected, DotsConfig)
		})
	}
}

func TestParse(t *testing.T) {
	abs, err := filepath.Abs(".")
	assert.NoErrorf(t, err, "failed to setup config_test.go testing couldn't get absolute path of `.`: %w", err)
	current, previous := abs, ""
	for current != previous {
		current, previous = filepath.Dir(current), current
	}
	_, err = config.Parse()
	assert.EqualError(t, err, fmt.Sprintf("failed to find `.dots.ya?ml` starting from: `%s` reached mount point `%s`", abs, current))
}
