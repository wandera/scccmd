package cmd

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var loglevel string

var rootCmd = &cobra.Command{
	Use:               "scccmd",
	DisableAutoGenTag: true,
	Short:             "Spring Cloud Config management tool",
	Long: `Commandline tool used for managing configuration from Spring Cloud Config Server.
Tool currently provides functionality to get (download) config file from server.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		lvl, err := log.ParseLevel(loglevel)
		if err != nil {
			return err
		}

		log.SetLevel(lvl)
		return nil
	},
	SilenceErrors: true,
}

func init() {
	rootCmd.PersistentFlags().StringVar(&loglevel, "log-level", "info", fmt.Sprintf("command log level (options: %s)", log.AllLevels))

	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(genDocCmd)
	rootCmd.AddCommand(encryptCmd)
	rootCmd.AddCommand(decryptCmd)
	rootCmd.AddCommand(webhookCmd)
	rootCmd.AddCommand(diffCmd)
	rootCmd.AddCommand(versionCmd)
}

//Execute run root command (main entrypoint)
func Execute() error {
	return rootCmd.Execute()
}
