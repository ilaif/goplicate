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

func TestRun_Error_SyncingToNonExistentFile(t *testing.T) {
	r := require.New(t)

	defer testutils.PrepareWorkdir(t, "../examples/sync-initial", ".")()

	cfg := &config.ProjectConfig{
		Targets: []config.Target{
			{
				Path:        "config.yaml",
				Source:      config.Source{Path: "./shared/config.yaml"},
				SyncInitial: false,
			},
		},
	}
	cloner := &mocks.ClonerMock{}
	opts := pkg.NewRunOpts(false, true, false, false, false, false, "")

	r.ErrorContains(pkg.Run(context.TODO(), cfg, cloner, opts), "Failed to read file")
}

func TestRun_Success_SyncingToNonExistentFile_WithSyncInitial(t *testing.T) {
	r := require.New(t)

	defer testutils.PrepareWorkdir(t, "../examples/sync-initial", ".")()

	cfg := &config.ProjectConfig{
		Targets: []config.Target{
			{
				Path:        "config.yaml",
				Source:      config.Source{Path: "./shared/config.yaml"},
				SyncInitial: true,
			},
		},
	}
	cloner := &mocks.ClonerMock{}
	opts := pkg.NewRunOpts(false, true, false, false, false, false, "")

	r.NoError(pkg.Run(context.TODO(), cfg, cloner, opts))

	testutils.RequireFileContains(r, "config.yaml", "key: value")
}
