package cmd_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ilaif/goplicate/pkg/cmd"
	"github.com/ilaif/goplicate/pkg/cmd/testutils"
)

func TestSyncCmd_Simple(t *testing.T) {
	r := require.New(t)

	defer testutils.PrepareWorkdir(t, "../../examples", "projects-simple")()

	testutils.RequireFileContains(r, "../simple/repo-1/.eslintrc.js", "indent: ['error', 4]")
	testutils.RequireFileContains(r, "../simple/repo-2/.eslintrc.js", "indent: ['error', 4]")

	syncCmd := cmd.NewSyncCmd()
	syncCmd.SetArgs([]string{"--confirm"})

	r.NoError(syncCmd.Execute())

	testutils.RequireFileContains(r, "../simple/repo-1/.eslintrc.js", "indent: ['error', 2]")
	testutils.RequireFileContains(r, "../simple/repo-2/.eslintrc.js", "indent: ['error', 2]")
}

func TestSyncCmd_RemoteGit(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test in short mode")
	}

	r := require.New(t)

	defer testutils.PrepareWorkdir(t, "../../examples", "projects-simple-remote-git")()

	testutils.RequireFileContains(r, "../simple-remote-git/repo-1/.eslintrc.js", "indent: ['error', 4]")
	testutils.RequireFileContains(r, "../simple-remote-git/repo-2/.eslintrc.js", "indent: ['error', 4]")

	syncCmd := cmd.NewSyncCmd()
	syncCmd.SetArgs([]string{"--confirm"})

	r.NoError(syncCmd.Execute())

	testutils.RequireFileContains(r, "../simple-remote-git/repo-1/.eslintrc.js", "indent: ['error', 2]")
	testutils.RequireFileContains(r, "../simple-remote-git/repo-2/.eslintrc.js", "indent: ['error', 2]")
}
