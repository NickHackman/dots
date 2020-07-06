package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

// DotsConfig a Dots config
type DotsConfig struct {
	Name     string    `yaml:"name"`     // Name that recognizes a set of dotfiles generally YourNameOrUsername/dotfiles
	License  string    `yaml:"license"`  // License used for dotfiles
	URL      string    `yaml:"URL"`      // URL to upstream
	Dotfiles []Dotfile `yaml:"dotfiles"` // Dotfiles themselves
}

// Dotfile a specific dotfile
type Dotfile struct {
	Name            string `yaml:"name"`             // Name that will be used to identify this specific Dotfile
	Description     string `yaml:"description"`      // Describe this specific dotfile or collection of dotfiles
	Source          string `yaml:"source"`           // Path to this dotfile
	Destination     string `yaml:"destination"`      // Path to install to
	InstallChildren bool   `yaml:"install_children"` // If true dictates that this dotfile is a logical organization of multiple dotfiles
}

const (
	configRegexp = `\.dots\.ya?ml`
)

// MountPointError is an error dictating that when searching for a `.dots.ya?ml` file
// one could not be found until it reached the mount point
type MountPointError struct {
	StartPoint string // Directory where search started at
	EndPoint   string // Highest directory found that ended discovery, generally `/` (on Unix)
}

// Error returns a String stating that the start p
func (mpe *MountPointError) Error() string {
	return fmt.Sprintf("failed to find `.dots.ya?ml` starting from: `%s` reached mount point `%s`", mpe.StartPoint, mpe.EndPoint)
}

// Finds .dots.ya?ml going upward till root is reached
//
// Root is determined if calling filepath.Dir(current) results in the same
// path as seen in this example https://golang.org/pkg/path/filepath/#Dir
func findConfig(startDir string) (string, error) {
	previous, current := "", startDir
	configRegex := regexp.MustCompile(configRegexp)

	for previous != current {
		files, err := ioutil.ReadDir(current)
		if err != nil {
			return "", fmt.Errorf("failed to read directory `%s`: %w", current, err)
		}
		for _, file := range files {
			if configRegex.MatchString(file.Name()) {
				return filepath.Join(current, file.Name()), nil
			}
		}
		previous, current = current, filepath.Dir(current)
	}
	return "", &MountPointError{StartPoint: startDir, EndPoint: current}
}

// Parse a `.dots.ya?ml`, starting from `start` progress upwards towards mount point.
// If mount point is reached a `MountPointError` is returned.
//
// Expects to find `.dots.ya?ml` in a parent dirctory of the current directory
func Parse(start string) (*DotsConfig, error) {
	abs, err := filepath.Abs(start)
	if err != nil {
		return nil, err
	}

	configPath, err := findConfig(abs)
	if err != nil {
		return nil, err
	}
	return ParseFile(configPath)
}

// Expands environment variables and '~' in Destination
//
// If Dotfile.Destination isn't set, set it to its default value `~/.config/$name`
func (dot *Dotfile) expandDestination() error {
	if dot.Destination == "" {
		configDir, err := os.UserConfigDir()
		if err != nil {
			return fmt.Errorf("failed to get User config directory: %w", err)
		}
		dot.Destination = filepath.Join(configDir, dot.Name)
	}
	dot.Destination = os.ExpandEnv(dot.Destination)

	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get current user home directory: %w", err)
	}

	if dot.Destination == "~" {
		dot.Destination = home
	} else if strings.HasPrefix(dot.Destination, "~/") {
		dot.Destination = filepath.Join(home, dot.Destination[2:])
	}
	return nil
}

// Expands environment variables and <root> in Source
//
// If Dotfile.Source isn't set, set it to its default value `<root>/$name`
func (dot *Dotfile) expandSource(projectRoot string) {
	if dot.Source == "" {
		dot.Source = fmt.Sprintf("<root>%c%s", os.PathSeparator, dot.Name)
	}
	dot.Source = os.ExpandEnv(dot.Source)
	if strings.HasPrefix(dot.Source, "<root>") {
		dot.Source = filepath.Join(projectRoot, dot.Source[6:])
	}
}

// ParseFile parses a 'dots.(yml|yaml)' file
func ParseFile(path string) (*DotsConfig, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file `%s`: %w", path, err)
	}
	dotsConf := DotsConfig{}
	if err := yaml.Unmarshal(bytes, &dotsConf); err != nil {
		return nil, fmt.Errorf("failed to parse `%s`: %w", path, err)
	}

	projectRoot := filepath.Dir(path)
	for i := range dotsConf.Dotfiles {
		dotsConf.Dotfiles[i].expandSource(projectRoot)
		if err := dotsConf.Dotfiles[i].expandDestination(); err != nil {
			return nil, err
		}
	}
	return &dotsConf, nil
}
