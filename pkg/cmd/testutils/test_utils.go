package testutils

import (
	"os"
	"testing"

	cp "github.com/otiai10/copy"
	"github.com/stretchr/testify/require"
)

func PrepareWorkdir(t *testing.T, source string) func() {
	r := require.New(t)

	dir, err := os.MkdirTemp(os.TempDir(), "_goplicate_"+t.Name())
	r.NoError(err)
	r.NoError(cp.Copy(source, dir))
	origWd, err := os.Getwd()
	r.NoError(err)
	r.NoError(os.Chdir(dir))

	return func() {
		os.RemoveAll(dir)
		_ = os.Chdir(origWd)
	}
}

func RequireFileContains(r *require.Assertions, filepath string, contains string) {
	bytes, err := os.ReadFile(filepath)
	r.NoError(err)
	contents := string(bytes)
	r.Contains(contents, contains)
}
