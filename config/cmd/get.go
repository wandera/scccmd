package cmd

import (
	"fmt"
	"github.com/WanderaOrg/scccmd/config/client"
	"github.com/spf13/cobra"
	"io/ioutil"
)

var (
	source       string
	application  string
	profile      string
	label        string
	format       string
	destination  string
	fileMappings FileMappings
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get the config from the given config server",
}

var getValuesCmd = &cobra.Command{
	Use:   "values",
	Short: "Get the config values in specified format from the given config server",
	RunE: func(cmd *cobra.Command, args []string) error {
		return executeGetValues(args)
	},
}

var getFilesCmd = &cobra.Command{
	Use:   "files",
	Short: "Get the config files from the given config server",
	RunE: func(cmd *cobra.Command, args []string) error {
		return executeGetFiles(args)
	},
}

func executeGetValues(args []string) error {
	ext, err := client.ParseExtension(format)

	if err != nil {
		return err
	}

	resp, err := client.
		NewClient(client.Config{URI: source, Profile: profile, Application: application, Label: label}).
		FetchAs(ext)

	if err != nil {
		return err
	}

	if destination != "" {
		if Verbose {
			fmt.Println("Config server response:")
			fmt.Println(resp)
		}

		if err = ioutil.WriteFile(destination, []byte(resp), 0644); err != nil {
			return err
		}

		if Verbose {
			fmt.Println("Response written to: ", destination)
		}
	} else {
		fmt.Print(resp)
	}

	return nil
}

func executeGetFiles(args []string) error {
	for _, mapping := range fileMappings.Mappings() {
		resp, err := client.
			NewClient(client.Config{URI: source, Profile: profile, Application: application, Label: label}).
			FetchFile(mapping.source)

		if err != nil {
			return err
		}

		if Verbose {
			fmt.Println("Config server response:")
			fmt.Println(string(resp))
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
	getCmd.AddCommand(getFilesCmd)
	getCmd.AddCommand(getValuesCmd)
	getCmd.PersistentFlags().StringVarP(&source, "source", "s", "", "address of the config server")
	getCmd.PersistentFlags().StringVarP(&application, "application", "a", "", "name of the application to get the config for")
	getCmd.PersistentFlags().StringVarP(&profile, "profile", "p", "default", "configuration profile")
	getCmd.PersistentFlags().StringVarP(&label, "label", "l", "master", "configuration label")
	getCmd.MarkFlagRequired("source")
	getCmd.MarkFlagRequired("application")

	getFilesCmd.Flags().VarP(&fileMappings, "files", "f", "files to get in form of source:destination pairs, example '--files application.yaml:config.yaml'")
	getFilesCmd.MarkFlagRequired("files")

	getValuesCmd.Flags().StringVarP(&format, "format", "f", "yaml", "output format might be one of 'json|yaml|properties'")
	getValuesCmd.Flags().StringVarP(&destination, "destination", "d", "", "destination file name")
}
