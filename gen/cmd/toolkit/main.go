package main

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/iotexproject/Bumblebee/gen/cmd/toolkit/cmd"
)

var command = &cobra.Command{
	Use:     "toolkit",
	Version: "0.0.1",
}

func init() {
	verbose := false
	command.PersistentFlags().BoolVarP(&verbose, "verbose", "v", verbose, "")
	command.AddCommand(cmd.Cmd)
	command.AddCommand(cmd.Patch)

}

func main() {
	if err := command.Execute(); err != nil {
		command.PrintErr(err)
		os.Exit(-1)
	}
}
