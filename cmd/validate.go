package cmd

import (
	"fmt"
	"os"

	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"

	"github.com/NickHackman/dots/config"
)

var configPath string

// validateCmd represents the validate command
var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate your .dots.ya?ml configuration file",
	Long: `Validate your .dots.ya?ml configuration file for common issues.

Validate will find the closest '.dots.ya?ml' file By starting at the current working directory
and progressing upwards until it finds a configuration file or the mount point ('/' on unix systems).

Use the '--config' or '-c' flag in order to pass a path to a dots configuration file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if configPath == "" {
			path, err := os.Getwd()
			if err != nil {
				fmt.Printf("%s: failed to get current working directory - %v", aurora.Red("Error"), err)
				return err
			}

			if configPath, err = config.FindConfig(path); err != nil {
				fmt.Printf("%s: %v\n", aurora.Red("Error"), err)
				os.Exit(1)
			}
		}

		err := config.Validate(configPath)
		if err == nil {
			return nil
		}

		for _, warn := range err.Warnings {
			fmt.Printf("%s: %s\n", aurora.Yellow("Warning"), warn.Message)
			if warn.Recommendation != "" {
				fmt.Printf("%s: %s\n", aurora.Blue("Info"), warn.Recommendation)
			}
		}

		if err.IsErr() {
			fmt.Printf("%s: %v\n", aurora.Red("Error"), err)
		}

		os.Exit(1)

		return nil // Unreachable
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)

	validateCmd.Flags().StringVarP(&configPath, "config", "c", "", "Path to `.dots.yml` file")
}
