package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	git "github.com/libgit2/git2go/v30"
	"gopkg.in/yaml.v3"
)

// A Dots config
type DotsConfig struct {
	Name     string    `yaml:"name"`     // Name that recognizes a set of dotfiles generally YourNameOrUsername/dotfiles
	License  string    `yaml:"license"`  // License used for dotfiles
	URL      string    `yaml:"URL"`      // URL to upstream
	Dotfiles []Dotfile `yaml:"dotfiles"` // Dotfiles themselves
}

// A specific dotfile
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

// Finds .dots.ya?ml in a directory
//
// In the case both .dots.yaml and .dots.yml exist .dots.yml will be chosen
func findConfigInDir(dir string) (string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return "", fmt.Errorf("failed to read directory `%s`: %w", dir, err)
	}

	configRegex := regexp.MustCompile(configRegexp)
	for _, file := range files {
		if configRegex.MatchString(file.Name()) {
			return file.Name(), nil
		}
	}
	return "", fmt.Errorf("'.dots.ya?ml' doesn't exist in dir %s", dir)
}

// Parse 'dots.(yml|yaml) in root of the git repository
//
// expects to find a 'dots.(yml|yaml)' file in the root of the git repository
func Parse() (*DotsConfig, error) {
	path, err := git.Discover(".", true, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to discover git repository in `.`: %w", err)
	}
	gitDir := filepath.Dir(path)
	root := filepath.Dir(gitDir)
	configPath, err := findConfigInDir(root)
	if err != nil {
		return nil, err
	}
	return ParseFile(configPath)
}

// Expands environemnt variables and '~' in Destination
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

// Parses 'dots.(yml|yaml)' file
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

// Parse 'dots.(yml|yaml) in root of the git repository
//
// expects to find a 'dots.(yml|yaml)' file in the root of the git repository
func ParseGit(gitRepo *git.Repository) (*DotsConfig, error) {
	dir := filepath.Dir(gitRepo.Workdir())
	configPath, err := findConfigInDir(dir)
	if err != nil {
		return nil, err
	}

	return ParseFile(configPath)
}