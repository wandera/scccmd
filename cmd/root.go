package cmd

import (
	"github.com/spf13/cobra"
)

//Verbose verbose logging turned on
var Verbose bool

var rootCmd = &cobra.Command{
	Use:               "scccmd",
	DisableAutoGenTag: true,
	Short:             "Spring Cloud Config management tool",
	Long: `Commandline tool used for managing configuration from Spring Cloud Config Server.
Tool currently provides functionality t get (download) config file from server.`,
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(genDocCmd)
	rootCmd.AddCommand(initializerCmd)
	rootCmd.AddCommand(encryptCmd)
	rootCmd.AddCommand(decryptCmd)
}

//Execute run root command (main entrypoint)
func Execute() error {
	return rootCmd.Execute()
}
