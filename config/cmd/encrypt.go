package cmd

import (
	"fmt"
	"github.com/WanderaOrg/scccmd/config/client"
	"github.com/spf13/cobra"
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
		return executeEncrypt(args)
	},
}

func executeEncrypt(args []string) error {
	if ep.value == "" {
		bytes, err := ioutil.ReadAll(os.Stdin)

		ep.value = string(bytes)
		if err != nil {
			return err
		}
	}

	if res, err := client.NewClient(client.Config{
		URI: ep.source,
	}).Encrypt(ep.value); err == nil {
		fmt.Println(res)
		return nil
	} else {
		return err
	}
}

func init() {
	encryptCmd.Flags().StringVarP(&ep.source, "source", "s", "", "address of the config server")
	encryptCmd.Flags().StringVar(&ep.value, "value", "", "value to encrypt *WARNING* unsafe use standard-in instead")
	encryptCmd.MarkFlagRequired("source")
}
