package cmd_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ilaif/goplicate/pkg/cmd"
	"github.com/ilaif/goplicate/pkg/cmd/testutils"
)

func TestRunCmd(t *testing.T) {
	r := require.New(t)

	defer testutils.PrepareWorkdir(t, "../../examples/simple", "repo-1")()

	testutils.RequireFileContains(r, ".eslintrc.js", "indent: ['error', 4]")

	runCmd := cmd.NewRunCmd()
	runCmd.SetArgs([]string{"--confirm"})

	r.NoError(runCmd.Execute())

	testutils.RequireFileContains(r, ".eslintrc.js", "indent: ['error', 2]")
}
