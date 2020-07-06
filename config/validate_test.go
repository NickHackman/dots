package config_test

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/NickHackman/dots/config"

	"github.com/stretchr/testify/assert"
)

func TestValidateValid(t *testing.T) {
	testData, err := pathToTestData()
	assert.NoErrorf(t, err, "failed to setup validate_test.go testing: %w", err)

	files, err := ioutil.ReadDir(testData)
	assert.NoErrorf(t, err, "failed to setup validate_test.go testing can't readir `%s`: %w", testData, err)

	for _, file := range files {
		if file.IsDir() || strings.HasPrefix(file.Name(), "invalid") {
			continue
		}

		t.Run(file.Name(), func(t *testing.T) {
			path := filepath.Join(testData, file.Name())
			assert.FileExists(t, path)
			validErr := config.Validate(path)
			if !assert.Nil(t, validErr) {
				t.Errorf("ValidationError was expected to be nil: %v\n", validErr.Err)
				for _, warn := range validErr.Warnings {
					t.Errorf("Message: %s Recommendation: %s\n", warn.Message, warn.Recommendation)
				}
			}
		})
	}
}

func TestValidateInvalid(t *testing.T) {
	testData, err := pathToTestData()
	assert.NoErrorf(t, err, "failed to setup validate_test.go testing: %w", err)

	homeDir, err := os.UserHomeDir()
	assert.NoErrorf(t, err, "failed to setup validate_test.go testing can't locate `$HOME`: %w", err)

	bspwmGitKeepPath := filepath.Join(testData, "bspwm", ".gitkeep")
	if _, err = os.Stat(bspwmGitKeepPath); !os.IsNotExist(err) {
		err = os.Remove(bspwmGitKeepPath)
		assert.NoErrorf(t, err, "failed to setup validate_test.go testing can't remove `%s`: %w", bspwmGitKeepPath, err)
	}

	tests := []struct {
		path            string
		validationError *config.ValidationError
	}{
		{
			path: "invalid-blank-name.yml",
			validationError: &config.ValidationError{
				Warnings: []*config.Warning{
					{
						Message:        "dots config name shouldn't be left blank, isn't directly installable",
						Recommendation: "set name to default value `YourName/dotfiles`",
					},
				},
			},
		},
		{
			path: "invalid-dot-blank-name.yml",
			validationError: &config.ValidationError{
				Err: errors.New("dotfile number `1` name is blank, but field is required"),
			},
		},
		{
			path: "invalid-blank-license.yml",
			validationError: &config.ValidationError{
				Err: errors.New("license is required, if you're not sure which license consult https://choosealicense.com/"),
			},
		},
		{
			path: "invalid-duplicate-dot-names.yml",
			validationError: &config.ValidationError{
				Err: errors.New("dotfiles with index `3` and `1` both have the same name `bspwm`"),
			},
		},
		{
			path: "invalid-duplicate-dot-destinations.yml",
			validationError: &config.ValidationError{
				Err: fmt.Errorf("dotfiles `bspwm` and `keybinds` have the same destination `%s` and will overwrite one another", homeDir),
			},
		},
		{
			path: "invalid-duplicate-dot-sources.yml",
			validationError: &config.ValidationError{
				Err: fmt.Errorf("dotfiles `bspwm` and `keybinds` have the same source `%s`", filepath.Join(testData, "bspwm")),
			},
		},
		{
			path: "invalid-duplicate-dot-descriptions.yml",
			validationError: &config.ValidationError{
				Warnings: []*config.Warning{{Message: "dotfiles bspwm and keybinds have the same description `description`"}},
				Err:      nil,
			},
		},
		{
			path: "invalid-dot-non-existing-source.yml",
			validationError: &config.ValidationError{
				Err: fmt.Errorf("dotfile `test2` source field `%s` does not exist", filepath.Join(testData, "test2")),
			},
		},
		{
			path: "invalid-dot-install-children-no-children.yml",
			validationError: &config.ValidationError{
				Err: fmt.Errorf("dotfile `bspwm` has `install_children` set, but has 0 children in source `%s`", filepath.Join(testData, "bspwm")),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.path, func(t *testing.T) {
			path := filepath.Join(testData, test.path)
			validationError := config.Validate(path)
			assert.Equal(t, validationError, test.validationError)
		})
	}

	_, err = os.Create(bspwmGitKeepPath)
	assert.NoError(t, err)
}
