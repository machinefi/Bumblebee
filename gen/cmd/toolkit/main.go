package main

import (
	"os"

	"github.com/saitofun/qkit/gen/cmd/toolkit/cmd"
	"github.com/spf13/cobra"
)

var command = &cobra.Command{
	Use:     "gen",
	Version: "0.0.1",
}

func init() {
	verbose := false
	command.PersistentFlags().BoolVarP(&verbose, "verbose", "v", verbose, "")
	command.AddCommand(cmd.Cmd)
}

func main() {
	if err := command.Execute(); err != nil {
		command.PrintErr(err)
		os.Exit(-1)
	}
}