package mocks

import (
	"context"

	"github.com/ilaif/goplicate/pkg/git"
)

type ClonerMock struct {
}

func (c *ClonerMock) Clone(
	ctx context.Context,
	uri, branch, fixedClonePath string,
) (clonePath string, err error) {
	return "", nil
}

func (c *ClonerMock) Close() {}

var _ git.Cloner = &ClonerMock{}
