package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/wandera/scccmd/pkg/client"
	"io/ioutil"
	"os"
)

var dp = struct {
	source string
	value  string
}{}

var decryptCmd = &cobra.Command{
	Use:   "decrypt",
	Short: "Decrypt the value server-side and prints the response",
	RunE: func(cmd *cobra.Command, args []string) error {
		return ExecuteDecrypt()
	},
}

// ExecuteDecrypt runs decrypt cmd
func ExecuteDecrypt() error {
	if dp.value == "" {
		bytes, err := ioutil.ReadAll(os.Stdin)

		dp.value = string(bytes)
		if err != nil {
			return err
		}
	}

	res, err := client.NewClient(client.Config{
		URI: dp.source,
	}).Decrypt(dp.value)

	if err == nil {
		fmt.Println(res)
	}

	return err
}

func init() {
	decryptCmd.Flags().StringVarP(&dp.source, "source", "s", "", "address of the config server")
	decryptCmd.Flags().StringVar(&dp.value, "value", "", "value to decrypt *WARNING* unsafe use standard-in instead")
	_ = decryptCmd.MarkFlagRequired("source")
}
