package cmd

import (
	"github.com/spf13/cobra"
)

var Verbose bool

var rootCmd = &cobra.Command{
	Use:               "config",
	DisableAutoGenTag: true,
	Short:             "Spring Cloud Config management tool",
	Long:
	`Commandline tool used for managing configuration from Spring Cloud Config Server.
Tool currently provides functionality t get (download) config file from server.`,
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(genDocCmd)
}

func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		return err
	}
	return nil
}
