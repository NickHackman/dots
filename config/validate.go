package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
)

// ValidationError a validation Error and or Warnings
type ValidationError struct {
	Warnings []*Warning // A slice of Warnings and their corresponding recommendations
	Err      error      // Error that occurred
}

// Error is shorthand for ValidationError.Err.Error(), empty string if no Error present
func (ve *ValidationError) Error() string {
	if ve.Err == nil {
		return ""
	}
	return ve.Err.Error()
}

// IsErr checks to see if ValidationError is an Error
func (ve *ValidationError) IsErr() bool {
	return ve.Err != nil
}

// Warning a message of what's wrong and possibly a recommendation on how to fix it
type Warning struct {
	Message        string // Warning message
	Recommendation string // Recommendation on how to fix this warning
}

type validator struct {
	dotsConf *DotsConfig
	validErr *ValidationError
}

// Validate validates a dots.yml file
func Validate(path string) *ValidationError {
	dotsConf, err := ParseFile(path)
	if err != nil {
		return &ValidationError{Err: err}
	}
	validator := &validator{dotsConf: dotsConf, validErr: &ValidationError{}}
	validator.validateName()

	if validator.dotsConf.License == "" {
		validator.validErr.Err = errors.New("license is required, if you're not sure which license consult https://choosealicense.com/")
		return validator.validErr
	}

	if err = validator.validateDots(); err != nil {
		validator.validErr.Err = err
		return validator.validErr
	}

	if len(validator.validErr.Warnings) != 0 {
		return validator.validErr
	}

	return nil
}

// Validates .dots.ya?ml Name to not be blank
func (v *validator) validateName() {
	if v.dotsConf.Name != "" {
		return
	}

	Message := "dots config name shouldn't be left blank, isn't directly installable"
	Recommendation := "set name to default value `YourName/dotfiles`"
	warn := &Warning{Message, Recommendation}
	v.validErr.Warnings = append(v.validErr.Warnings, warn)
}

// Validate dotfiles
//
// Check individual fields
//
// Name            - MUST exist and not be empty
// Source          - MUST exist
// Destination     - Probably shouldn't equal `~` or `/`
// Description     - shouldn't be empty
// InstallChildren - MUST have children
func (v *validator) validateDots() error {
	if v.dotsConf.Dotfiles == nil || len(v.dotsConf.Dotfiles) == 0 {
		return nil
	}

	for i, dot := range v.dotsConf.Dotfiles {
		if dot.Name == "" {
			return fmt.Errorf("dotfile number `%d` name is blank, but field is required", i+1)
		}

		if _, err := os.Stat(dot.Source); os.IsNotExist(err) {
			return fmt.Errorf("dotfile `%s` source field `%s` does not exist", dot.Name, dot.Source)
		}

		if dot.Description == "" {
			Message := fmt.Sprintf("dotfile `%s` description shouldn't be left blank", dot.Name)
			v.validErr.Warnings = append(v.validErr.Warnings, &Warning{Message: Message})
		}

		if dot.InstallChildren {
			files, err := ioutil.ReadDir(dot.Source)
			if err != nil {
				return fmt.Errorf("dotfile `%s` failed to read directory `%s`: %w", dot.Name, dot.Source, err)
			}
			if len(files) == 0 {
				return fmt.Errorf("dotfile `%s` has `install_children` set, but has 0 children in source `%s`", dot.Name, dot.Source)
			}
		}
	}
	fields := reflect.TypeOf(v.dotsConf.Dotfiles[0])
	numFields := fields.NumField()
	for i := 0; i < numFields; i++ {
		if err := v.validateDuplicateDotVals(i); err != nil {
			return err
		}
	}
	return nil
}

// Validates duplicate fields in Dots by fieldName, hard errors on duplicate names, destinations, and sources
func (v *validator) validateDuplicateDotVals(fieldIndex int) error {
	var DupMap = make(map[string]struct {
		name  string
		index int
	})

	for i, dot := range v.dotsConf.Dotfiles {
		fieldValue := reflect.ValueOf(dot).Field(fieldIndex).String()
		fieldName := reflect.TypeOf(dot).Field(fieldIndex).Name
		if prevDot, ok := DupMap[fieldValue]; ok {
			var Message string
			switch fieldName {
			case "Name":
				return fmt.Errorf("dotfiles with index `%d` and `%d` both have the same name `%s`", i+1, prevDot.index+1, dot.Name)
			case "Destination":
				return fmt.Errorf("dotfiles `%s` and `%s` have the same destination `%s` and will overwrite one another", prevDot.name, dot.Name, dot.Destination)
			case "Source":
				return fmt.Errorf("dotfiles `%s` and `%s` have the same source `%s`", prevDot.name, dot.Name, dot.Source)
			case "Description":
				Message = fmt.Sprintf("dotfiles %s and %s have the same description `%s`", prevDot.name, dot.Name, dot.Description)
			case "InstallChildren":
				continue
			default:
				panic(fmt.Sprintf("Unknown field `%s` in Dotfile if duplicates matters please implement a case for it in validateDuplicateDotVals; otherwise, exclude it.", fieldName))
			}
			warn := &Warning{Message: Message}
			v.validErr.Warnings = append(v.validErr.Warnings, warn)
			continue
		}
		DupMap[fieldValue] = struct {
			name  string
			index int
		}{name: dot.Name, index: i}
	}

	return nil
}
