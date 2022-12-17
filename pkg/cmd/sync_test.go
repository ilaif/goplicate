package cmd_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ilaif/goplicate/pkg/cmd"
	"github.com/ilaif/goplicate/pkg/cmd/testutils"
	"github.com/ilaif/goplicate/pkg/utils"
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

	rootCmd := cmd.NewRootCmd("")
	rootCmd.SetArgs([]string{"sync", "--confirm", "--disable-cleanup", "--debug"})

	r.NoError(rootCmd.Execute())

	fmt.Println(utils.MustGetwd())

	testutils.RequireFileContains(r,
		"../cloned/repo-1/examples/simple-remote-git/repo-1/.eslintrc.js",
		"indent: ['error', 2]",
	)
	testutils.RequireFileContains(r,
		"../cloned/repo-2/examples/simple-remote-git/repo-2/.eslintrc.js",
		"indent: ['error', 2]",
	)
}
