package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"github.com/wanderaorg/scccmd/config/client"
)

var (
	source       string
	application  string
	profile      string
	label        string
	fileMappings FileMappings
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the config from the given config server",
	RunE: func(cmd *cobra.Command, args []string) error {
		return executeInit(args)
	},
}

func executeInit(args []string) error {
	for _, mapping := range fileMappings.Mappings() {
		resp, err := client.
			NewClient(client.Config{URI: source, Profile: profile, Application: application, Label: label}).
			FetchFile(mapping.source)

		if err != nil {
			return err
		}

		if Verbose {
			fmt.Println("Config server response:")
			fmt.Println(resp)
		}

		if err = ioutil.WriteFile(mapping.destination, resp, 0644); err != nil {
			return err
		}

		if Verbose {
			fmt.Println("Response written to: ", mapping.destination)
		}
	}
	return nil
}

func init() {
	initCmd.Flags().StringVarP(&source, "source", "s", "", "address of the config server")
	initCmd.Flags().StringVarP(&application, "application", "a", "", "name of the application to get the config for")
	initCmd.Flags().StringVarP(&profile, "profile", "p", "default", "configuration profile")
	initCmd.Flags().StringVarP(&label, "label", "l", "master", "configuration label")
	initCmd.Flags().VarP(&fileMappings, "files", "f", "files to get in form of source:destination pairs, example '--files application.yaml:config.yaml'")
	initCmd.MarkFlagRequired("source")
	initCmd.MarkFlagRequired("application")
	initCmd.MarkFlagRequired("files")
}
