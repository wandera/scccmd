package cmd

import "testing"

func TestRootCommand(t *testing.T) {
	err := rootCmd.Execute()
	if err != nil {
		t.Error("Running root command should not throw exception", err)
	}
}
