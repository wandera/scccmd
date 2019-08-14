package cmd

import (
	"fmt"
	"github.com/WanderaOrg/scccmd/pkg/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
)

const stdoutPlaceholder = "-"

var gp = struct {
	source       string
	application  string
	profile      string
	label        string
	format       string
	destination  string
	fileMappings FileMappings
}{}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get the config from the given config server",
}

var getValuesCmd = &cobra.Command{
	Use:   "values",
	Short: "Get the config values in specified format from the given config server",
	RunE: func(cmd *cobra.Command, args []string) error {
		return ExecuteGetValues()
	},
}

var getFilesCmd = &cobra.Command{
	Use:   "files",
	Short: "Get the config files from the given config server",
	RunE: func(cmd *cobra.Command, args []string) error {
		return ExecuteGetFiles()
	},
}

// ExecuteGetValues runs get values cmd
func ExecuteGetValues() error {
	ext, err := client.ParseExtension(gp.format)

	if err != nil {
		return err
	}

	resp, err := client.
		NewClient(client.Config{URI: gp.source, Profile: gp.profile, Application: gp.application, Label: gp.label}).
		FetchAs(ext)

	if err != nil {
		return err
	}

	if gp.destination != "" {
		log.Debug("Config server response:")
		log.Debug(resp)

		if err = ioutil.WriteFile(gp.destination, []byte(resp), 0644); err != nil {
			return err
		}

		log.Debug("Response written to: ", gp.destination)

	} else {
		fmt.Print(resp)
	}

	return nil
}

// ExecuteGetFiles runs get files cmd
func ExecuteGetFiles() error {
	for _, mapping := range gp.fileMappings.Mappings() {
		resp, err := client.
			NewClient(client.Config{URI: gp.source, Profile: gp.profile, Application: gp.application, Label: gp.label}).
			FetchFileE(mapping.source)

		if err != nil {
			return err
		}

		log.Debug("Config server response:")
		log.Debug(string(resp))

		if mapping.destination == stdoutPlaceholder {
			_, _ = os.Stdout.Write(resp)
			fmt.Println()
			log.Debug("Response written to stdout")
		} else {
			if err = ioutil.WriteFile(mapping.destination, resp, 0644); err != nil {
				return err
			}

			log.Debug("Response written to: ", mapping.destination)
		}
	}
	return nil
}

func init() {
	getCmd.AddCommand(getFilesCmd)
	getCmd.AddCommand(getValuesCmd)
	getCmd.PersistentFlags().StringVarP(&gp.source, "source", "s", "", "address of the config server")
	getCmd.PersistentFlags().StringVarP(&gp.application, "application", "a", "", "name of the application to get the config for")
	getCmd.PersistentFlags().StringVarP(&gp.profile, "profile", "p", "default", "configuration profile")
	getCmd.PersistentFlags().StringVarP(&gp.label, "label", "l", "master", "configuration label")
	_ = getCmd.MarkPersistentFlagRequired("source")
	_ = getCmd.MarkPersistentFlagRequired("application")

	getFilesCmd.Flags().VarP(&gp.fileMappings, "files", "f", "files to get in form of source:destination pairs, you can use - as a output to stdout, example '--files application.yaml:config.yaml'")
	_ = getFilesCmd.MarkFlagRequired("files")

	getValuesCmd.Flags().StringVarP(&gp.format, "format", "f", "yaml", "output format might be one of 'json|yaml|properties'")
	getValuesCmd.Flags().StringVarP(&gp.destination, "destination", "d", "", "destination file name")
}
