package cmd

import (
	"github.com/iotexproject/Bumblebee/kit/statusxgen"
	"github.com/iotexproject/Bumblebee/x/pkgx"
	"github.com/spf13/cobra"
)

func init() {
	cmd := &cobra.Command{
		Use:     "status",
		Aliases: []string{"status-error", "error"},
		Short:   "generate interfaces of status error",
		Run: func(cmd *cobra.Command, args []string) {
			run("status", func(pkg *pkgx.Pkg) Generator {
				g := statusxgen.New(pkg)
				g.Scan(args...)
				return g
			}, args...)
		},
	}

	Cmd.AddCommand(cmd)
}
