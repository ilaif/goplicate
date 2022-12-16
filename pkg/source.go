package pkg

import (
	"context"
	"path"

	"github.com/caarlos0/log"

	"github.com/ilaif/goplicate/pkg/git"
)

// ResolveSourcePath given a source, resolves it by cloning the repository (if applicable)
// and returning the directory of the source.
func ResolveSourcePath(ctx context.Context, source Source, workdir string, cloner *git.Cloner) (string, error) {
	log.Debugf("Resolving path of source '%s'...", source.String())

	var err error

	branch := source.Branch
	if source.Tag != "" {
		branch = source.Tag
	}

	dir := workdir
	if source.Repository != "" {
		dir, err = cloner.Clone(ctx, string(source.Repository), branch)
		if err != nil {
			return "", err
		}
	}

	return path.Join(dir, source.Path), nil
}
