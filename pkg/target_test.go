package pkg_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ilaif/goplicate/pkg"
	"github.com/ilaif/goplicate/pkg/cmd/testutils"
	"github.com/ilaif/goplicate/pkg/config"
	"github.com/ilaif/goplicate/pkg/mocks"
)

func TestRunTarget_Error_SyncingToNonExistentFile(t *testing.T) {
	r := require.New(t)

	defer testutils.PrepareWorkdir(t, "../examples/sync-initial", ".")()

	target := config.Target{
		Path:        "config.yaml",
		Source:      config.Source{Path: "./shared/config.yaml"},
		SyncInitial: false,
	}
	cloner := &mocks.ClonerMock{}

	_, err := pkg.RunTarget(context.TODO(), target, cloner, false, true)
	r.ErrorContains(err, "Failed to read file")
}

func TestRunTarget_Success_SyncingToNonExistentFile_WithSyncInitial(t *testing.T) {
	r := require.New(t)

	defer testutils.PrepareWorkdir(t, "../examples/sync-initial", ".")()

	target := config.Target{
		Path:        "config.yaml",
		Source:      config.Source{Path: "./shared/config.yaml"},
		SyncInitial: true,
	}
	cloner := &mocks.ClonerMock{}

	_, err := pkg.RunTarget(context.TODO(), target, cloner, false, true)
	r.NoError(err)

	testutils.RequireFileContains(r, "config.yaml", "key: value")
}
