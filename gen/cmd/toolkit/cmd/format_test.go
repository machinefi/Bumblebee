package cmd_test

import (
	"testing"

	"golang.org/x/mod/modfile"

	"github.com/machinefi/Bumblebee/gen/cmd/toolkit/cmd"
)

func TestModulePath(t *testing.T) {
	mod := modfile.ModulePath([]byte(`
module github.com/machinefi/Bumblebee

go 1.18
`))
	t.Logf(mod)

	mod = modfile.ModulePath([]byte(`xxx`))
	t.Logf(mod)
}

func TestPrepareArg(t *testing.T) {
	cmd.PrepareArgs()
}

func TestFormatRoot(t *testing.T) {
	_ = cmd.FormatRoot("github.com/machinefi/Bumblebee", "/Users/sincos/sincos/src/github.com/machinefi/Bumblebee/base/types/snowflake_id", true)
}
