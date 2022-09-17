package cmd

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ilaif/goplicate/pkg/cmd/testutils"
)

func TestSyncCmd(t *testing.T) {
	r := require.New(t)

	defer testutils.PrepareWorkdir(t, "testdata")()
	r.NoError(os.Chdir("./project-simple-valid"))

	testutils.RequireFileContains(r, "../simple-valid/.eslintrc.js", "indent: ['error', 4]")

	cmd := NewSyncCmd()
	cmd.SetArgs([]string{"--confirm"})

	r.NoError(cmd.Execute())

	testutils.RequireFileContains(r, "../simple-valid/.eslintrc.js", "indent: ['error', 2]")
}
