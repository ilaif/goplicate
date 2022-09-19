package cmd

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ilaif/goplicate/pkg/cmd/testutils"
)

func TestRunCmd(t *testing.T) {
	r := require.New(t)

	defer testutils.PrepareWorkdir(t, "../../examples/simple-valid", "project-1")()

	testutils.RequireFileContains(r, ".eslintrc.js", "indent: ['error', 4]")

	cmd := NewRunCmd()
	cmd.SetArgs([]string{"--confirm"})

	r.NoError(cmd.Execute())

	testutils.RequireFileContains(r, ".eslintrc.js", "indent: ['error', 2]")
}
