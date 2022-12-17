package testutils

import (
	"os"
	"path"
	"testing"

	cp "github.com/otiai10/copy"
	"github.com/stretchr/testify/require"

	"github.com/ilaif/goplicate/pkg/utils"
)

func PrepareWorkdir(t *testing.T, source string, cd string) func() {
	r := require.New(t)

	dir, err := os.MkdirTemp(os.TempDir(), "_goplicate_"+t.Name())
	r.NoError(err)
	r.NoError(cp.Copy(source, dir))
	origWd, err := os.Getwd()
	r.NoError(err)
	r.NoError(os.Chdir(path.Join(dir, cd)))

	return func() {
		os.RemoveAll(dir)
		_ = os.Chdir(origWd)
	}
}

func RequireFileContains(r *require.Assertions, filepath string, contains string) {
	bytes, err := os.ReadFile(path.Join(utils.MustGetwd(), filepath))
	r.NoError(err)
	contents := string(bytes)
	r.Contains(contents, contains)
}
