package cmd

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ilaif/goplicate/pkg/cmd/testutils"
)

func TestSyncCmd(t *testing.T) {
	r := require.New(t)

	defer testutils.PrepareWorkdir(t, "../../examples", "projects-simple")()

	testutils.RequireFileContains(r, "../simple/repo-1/.eslintrc.js", "indent: ['error', 4]")
	testutils.RequireFileContains(r, "../simple/repo-2/.eslintrc.js", "indent: ['error', 4]")

	cmd := NewSyncCmd()
	cmd.SetArgs([]string{"--confirm"})

	r.NoError(cmd.Execute())

	testutils.RequireFileContains(r, "../simple/repo-1/.eslintrc.js", "indent: ['error', 2]")
	testutils.RequireFileContains(r, "../simple/repo-2/.eslintrc.js", "indent: ['error', 2]")
}
