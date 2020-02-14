package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/wandera/scccmd/pkg/client"
	"io/ioutil"
	"os"
)

var ep = struct {
	source string
	value  string
}{}

var encryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "Encrypt the value server-side and prints the response",
	RunE: func(cmd *cobra.Command, args []string) error {
		return ExecuteEncrypt()
	},
}

// ExecuteEncrypt runs encrypt cmd
func ExecuteEncrypt() error {
	if ep.value == "" {
		bytes, err := ioutil.ReadAll(os.Stdin)

		ep.value = string(bytes)
		if err != nil {
			return err
		}
	}

	res, err := client.NewClient(client.Config{
		URI: ep.source,
	}).Encrypt(ep.value)

	if err == nil {
		fmt.Println(res)
	}

	return err
}

func init() {
	encryptCmd.Flags().StringVarP(&ep.source, "source", "s", "", "address of the config server")
	encryptCmd.Flags().StringVar(&ep.value, "value", "", "value to encrypt *WARNING* unsafe use standard-in instead")
	_ = encryptCmd.MarkFlagRequired("source")
}
