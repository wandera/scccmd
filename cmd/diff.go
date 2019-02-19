package cmd

import (
	"fmt"
	"github.com/WanderaOrg/scccmd/pkg/client"
	"github.com/pmezard/go-difflib/difflib"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"net/http"
	"os"
	"strings"
)

var diffp = struct {
	source        string
	application   string
	profile       string
	label         string
	format        string
	destination   string
	files         string
	targetProfile string
	targetLabel   string
}{}

var diffCmd = &cobra.Command{
	Use:               "diff",
	Short:             "Diff the config from the given config server",
	PersistentPreRunE: validateDiffParams,
}

var diffValuesCmd = &cobra.Command{
	Use:   "values",
	Short: "Diff the config values in specified format from the given config server",
	RunE: func(cmd *cobra.Command, args []string) error {
		return ExecuteDiffValues()
	},
}

var diffFilesCmd = &cobra.Command{
	Use:   "files",
	Short: "Diff the config files from the given config server",
	RunE: func(cmd *cobra.Command, args []string) error {
		return ExecuteDiffFiles()
	},
}

func validateDiffParams(cmd *cobra.Command, args []string) error {
	err := rootCmd.PersistentPreRunE(cmd, args)
	if err != nil {
		return err
	}

	if diffp.targetProfile == "" {
		diffp.targetProfile = diffp.profile
	}
	return nil
}

//ExecuteDiffValues runs diff values cmd
func ExecuteDiffValues() error {
	ext, err := client.ParseExtension(diffp.format)

	if err != nil {
		return err
	}

	respA, err := client.
		NewClient(client.Config{URI: diffp.source, Profile: diffp.profile, Application: diffp.application, Label: diffp.label}).
		FetchAs(ext)

	if err != nil {
		return err
	}

	log.Debugf("Config server response for label %s, profile %s:", diffp.label, diffp.profile)
	log.Debug(string(respA))

	respB, err := client.
		NewClient(client.Config{URI: diffp.source, Profile: diffp.targetProfile, Application: diffp.application, Label: diffp.targetLabel}).
		FetchAs(ext)

	if err != nil {
		return err
	}

	log.Debugf("Config server response for label %s, profile %s:", diffp.targetLabel, diffp.targetProfile)
	log.Debug(string(respB))

	d := difflib.UnifiedDiff{
		A:       difflib.SplitLines(string(respA)),
		B:       difflib.SplitLines(string(respB)),
		Context: 3,
	}

	return difflib.WriteUnifiedDiff(os.Stdout, d)
}

//ExecuteDiffFiles runs diff files cmd
func ExecuteDiffFiles() error {
	errorHandler := func(data []byte, err error) []byte {
		if e, ok := err.(client.HTTPError); ok && e.StatusCode() == http.StatusNotFound {
			return []byte{}
		}
		fmt.Println(err.Error())
		return nil
	}

	for _, filename := range strings.Split(diffp.files, ",") {
		respA := client.
			NewClient(client.Config{URI: diffp.source, Profile: diffp.profile, Application: diffp.application, Label: diffp.label}).
			FetchFile(filename, errorHandler)

		if respA == nil {
			return fmt.Errorf("file %s for label %s and profile %s cannot be retrieved from remote server %s",
				filename, diffp.label, diffp.profile, diffp.source)
		}

		log.Debugf("Config server response for label %s, profile %s:", diffp.label, diffp.profile)
		log.Debug(string(respA))

		respB := client.
			NewClient(client.Config{URI: diffp.source, Profile: diffp.targetProfile, Application: diffp.application, Label: diffp.targetLabel}).
			FetchFile(filename, errorHandler)

		if respB == nil {
			return fmt.Errorf("file %s for label %s and profile %s cannot be retrieved from remote server %s",
				filename, diffp.targetLabel, diffp.targetProfile, diffp.source)
		}

		log.Debugf("Config server response for label %s, profile %s:", diffp.targetLabel, diffp.targetProfile)
		log.Debug(string(respB))

		d := difflib.UnifiedDiff{
			A:       difflib.SplitLines(string(respA)),
			B:       difflib.SplitLines(string(respB)),
			Context: 3,
		}

		diffString, err := difflib.GetUnifiedDiffString(d)

		if err != nil {
			return err
		}

		printFileDiff(diffString, filename)
	}
	log.Debug("Diff of files written to stdout")
	return nil
}

func printFileDiff(diffString string, filename string) {
	if len(diffString) > 0 {
		fmt.Printf("diff a/%s b/%s\n", filename, filename)
		fmt.Printf("--- a/%s profile=%s label=%s\n", filename, diffp.profile, diffp.label)
		fmt.Printf("+++ b/%s profile=%s label=%s\n", filename, diffp.targetProfile, diffp.targetLabel)
		fmt.Print(diffString)
	}
}

func init() {
	diffCmd.AddCommand(diffFilesCmd)
	diffCmd.AddCommand(diffValuesCmd)
	diffCmd.PersistentFlags().StringVarP(&diffp.source, "source", "s", "", "address of the config server")
	diffCmd.PersistentFlags().StringVarP(&diffp.application, "application", "a", "", "name of the application to get the config for")
	diffCmd.PersistentFlags().StringVar(&diffp.profile, "profile", "default", "configuration profile")
	diffCmd.PersistentFlags().StringVar(&diffp.label, "label", "master", "configuration label")
	diffCmd.PersistentFlags().StringVar(&diffp.targetLabel, "target-label", "", "second label to diff with")
	diffCmd.PersistentFlags().StringVar(&diffp.targetProfile, "target-profile", "", "second profile to diff with, --profile value will be used, if not defined")
	_ = diffCmd.MarkPersistentFlagRequired("source")
	_ = diffCmd.MarkPersistentFlagRequired("application")
	_ = diffCmd.MarkPersistentFlagRequired("target-label")

	diffFilesCmd.Flags().StringVarP(&diffp.files, "files", "f", "", "files to get in form of file1,file2, example '--files application.yaml,config.yaml'")
	_ = diffFilesCmd.MarkFlagRequired("files")

	diffValuesCmd.Flags().StringVarP(&diffp.format, "format", "f", "yaml", "output format might be one of 'json|yaml|properties'")
}
