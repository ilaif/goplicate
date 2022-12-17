package pkg

import (
	"context"
	"path"

	"github.com/caarlos0/log"
	"github.com/pkg/errors"

	"github.com/ilaif/goplicate/pkg/config"
	"github.com/ilaif/goplicate/pkg/git"
)

// ResolveSourcePath given a source, resolves it by cloning the repository (if applicable)
// and returning the directory of the source.
func ResolveSourcePath(ctx context.Context, source config.Source, workdir string, cloner git.Cloner) (string, error) {
	log.Debugf("Resolving path of source '%s'", source.String())

	var err error

	branch := source.Branch
	if source.Tag != "" {
		branch = source.Tag
	}

	dir := workdir
	if source.Repository != "" {
		absClonePath := ""
		if source.ClonePath != "" {
			absClonePath = path.Join(workdir, source.ClonePath)
		}
		dir, err = cloner.Clone(ctx, string(source.Repository), branch, absClonePath)
		if err != nil {
			return "", errors.Wrap(err, "Failed to clone repository")
		}
	}

	return path.Join(dir, source.Path), nil
}
