package pkg_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ilaif/goplicate/pkg"
	"github.com/ilaif/goplicate/pkg/cmd/testutils"
	"github.com/ilaif/goplicate/pkg/mocks"
	"github.com/ilaif/goplicate/pkg/shared"
)

func TestRun_Success_SyncConfig(t *testing.T) {
	r := require.New(t)

	defer testutils.PrepareWorkdir(t, "../examples/sync-config", ".")()

	cloner := &mocks.ClonerMock{}
	opts := pkg.NewRunOpts(false, true, false, false, false, false, "", "")

	sharedState := &shared.State{
		Message: "",
	}

	r.NoError(pkg.Run(context.TODO(), cloner, sharedState, opts))

	testutils.RequireFileContains(r, ".goplicate.yaml", "path: new.yaml")
	testutils.RequireFileContains(r, "new.yaml", "newKey: newValue")
}
