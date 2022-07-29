package cmd

import (
	"github.com/spf13/cobra"

	"github.com/iotexproject/Bumblebee/kit/enumgen"
	"github.com/iotexproject/Bumblebee/x/pkgx"
)

func init() {
	Cmd.AddCommand(&cobra.Command{
		Use:   "enum",
		Short: "generate interfaces of enumeration",
		Run: func(cmd *cobra.Command, args []string) {
			run("enum", func(pkg *pkgx.Pkg) Generator {
				g := enumgen.New(pkg)
				g.Scan(args...)
				return g
			}, args...)
		},
	})
}
