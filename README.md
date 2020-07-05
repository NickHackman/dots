# Dots

![build](https://github.com/NickHackman/dots/workflows/build/badge.svg)
![lint](https://github.com/NickHackman/dots/workflows/lint/badge.svg)
[![Coverage Status](https://coveralls.io/repos/github/NickHackman/dots/badge.svg?branch=master)](https://coveralls.io/github/NickHackman/dots?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/NickHackman/dots)](https://goreportcard.com/report/github.com/NickHackman/dots)

A distributed dotfile manager written in Go

All you have to do is create a `.dots.yml` file in the root of your git repository
and then it will be easily installable! For an example `.dots.yml` refer to [example](#example)

## Example

A `.dots.yml` configuration file looks like

```yaml
# This file is metadata that Dots will read to install
# the dotfiles in your repository locally on a new machine.
#
# This file is EXPECTED to be in the root of your git repository
# and be named `.dots.yml` or `.dots.yaml`. If it is not it's possible to set
# the configuration file in dots by passing the `--config=/path/to/config` flag,
# in this case please document this in your repository's README.
#
# Make sure to run
#
# $ dots validate
#
# to ensure that your configuration file is valid.

# Name of your dotfiles repository
#
# In the majority of cases should be `YourName/(dotfiles|dots|config)`
#
# Optional field
name: YourName/dotfiles

# The initialism of your dotfiles are licensed under
#
# Required field
license: GPLv3

# URL to your repository or upstream URL
#
# Required field
URL: https://github.com/NickHackman/dots

# List of dotfiles that will be installable and their required metadata
dotfiles:
  # Name of application to install
  #
  # This name MUST be unique and can be passed directly to Dots install.
  #
  # $ dots install bspwm
  - name: bspwm

    # Description of the application and your specific configuration
    #
    # Optional field
    description: A simple configuration file for the Binary Space Partition Window Manager

    # Where this dotfile is located in your repository.
    #
    # The default value is <root>/$name, meaning the root of this git
    # repository and the name of the current dotfile.
    #
    # This value CANNOT be outside of the git repository.
    #
    # To be platform agnostic write paths as if they were Unix (using `/` as the separator)
    # these will be resolved properly.
    #
    # Optional field
    source: <root>/bspwm

    # Where this dotfile should be installed to on a machine
    #
    # The default value is XDG_CONFIG_HOME/$name, which is generally `~/.config`
    # and the name of the current dotfile. Environment variables will be expanded.
    #
    # To be platform agnostic write paths as if they were Unix (using `/` as the separator)
    # these will be resolved properly.
    #
    # Optional field
    destination: ~/.config/bspwm

  - name: keybinds
    description: Keybindings that escape <-> capslock and handle function keys
    # In the case of a singular `~` it must be in either double or single quotes
    destination: "~"

    # $source is implicit as <root>/keybinds

    # Sometimes it makes sense to organize configuration files logically by directory.
    # In this case, this is effectively a shorthand for multiple dotfiles that will
    # be installed to $destination/$name.
    #
    # For instance
    #
    # keybinds/
    # |-- .xbindkeysrc
    # |-- .speedswapper
    #
    # $source MUST have children
    #
    # will be installed to
    #
    # ~/
    # |-- .xbindkeysrc
    # |-- .speedswapper
    #
    # Optional field
    install_children: true
```
